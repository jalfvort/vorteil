package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vorteil/vorteil/pkg/vpkg"
)

/**
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2020 vorteil.io Pty Ltd
 **/ //

var repositoriesCmd = &cobra.Command{
	Use:   "repositories",
	Short: "Interact with vorteil repositories",
}

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Create, List and Delete keys for authentication with Vorteil Repositories",
}

var createKeyCmd = &cobra.Command{
	Use:   "create NAME TOKEN",
	Short: "Create keys saves a file with the token to be referenced using a name",
	Args:  cobra.MaximumNArgs(2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("must provide a NAME and TOKEN")
		}
		usr, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		// Check if key storage exists if not create it
		pathCheck := filepath.Join(usr, ".vorteil", "repository-keys")
		_, err = os.Stat(pathCheck)
		if err != nil {
			err = os.MkdirAll(pathCheck, os.ModePerm)
			if err != nil {
				return err
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		key := args[1]

		usr, err := os.UserHomeDir()
		if err != nil {
			SetError(err, 1)
			return
		}
		pathCheck := filepath.Join(usr, ".vorteil", "repository-keys")

		// Check if file exists
		if !flagForce {
			// if stat returns no error return error saying you need to provide the force flag
			_, err := os.Stat(filepath.Join(pathCheck, name))
			if err == nil {
				SetError(errors.New("key file already exists provide --force to overwrite"), 2)
				return
			}
		}

		// Write key to a file under that keys directory
		err = ioutil.WriteFile(filepath.Join(pathCheck, name), []byte(key), os.ModePerm)
		if err != nil {
			SetError(err, 3)
			return
		}

		// Check if default and write another file called default under repository-keys
		if flagDefault {
			err = ioutil.WriteFile(filepath.Join(pathCheck, "default"), []byte(key), os.ModePerm)
			if err != nil {
				SetError(err, 4)
				return
			}
		}
	},
}

func init() {
	f := createKeyCmd.Flags()
	f.BoolVar(&flagDefault, "default", false, "save this key to use as default")
	f.BoolVar(&flagForce, "force", false, "force overwrite of key file")
}

var listKeysCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keys currently stored",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := os.UserHomeDir()
		if err != nil {
			SetError(err, 1)
			return
		}
		pathCheck := filepath.Join(usr, ".vorteil", "repository-keys")

		if _, err = os.Stat(pathCheck); os.IsNotExist(err) {
			fmt.Printf("no keys found\n")
			return
		}

		fis, err := ioutil.ReadDir(pathCheck)
		if err != nil {
			SetError(err, 2)
			return
		}

		for _, fi := range fis {
			if fi.Name() != "default" {
				fmt.Println(fi.Name())
			}
		}
	},
}

var deleteKeyCmd = &cobra.Command{
	Use:   "delete NAME",
	Short: "Delete a key currently stored",
	Args:  cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Must provide the name of the key you want to delete")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		usr, err := os.UserHomeDir()
		if err != nil {
			SetError(err, 1)
			return
		}
		pathCheck := filepath.Join(usr, ".vorteil", "repository-keys")
		path := filepath.Join(pathCheck, name)
		dpath := filepath.Join(pathCheck, "default")

		// before removing we should check if default is the same and delete that
		f1, err := ioutil.ReadFile(path)
		if err != nil {
			SetError(fmt.Errorf("%s keyfile does not exist", name), 2)
			return
		}

		f2, err := ioutil.ReadFile(dpath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				SetError(err, 3)
				return
			}
		}

		// If same bytes remove default aswell
		if bytes.Equal(f1, f2) {
			err = os.Remove(dpath)
			if err != nil {
				SetError(err, 4)
				return
			}
		}

		// Else just remove the keyfile
		err = os.Remove(filepath.Join(pathCheck, name))
		if err != nil {
			SetError(err, 5)
			return
		}

	},
}

var pushCmd = &cobra.Command{
	Use:   "push REPOSITORY ORG/BUCKET/APP SOURCE",
	Short: "Push to a repository",
	Long:  `The push command is a function for quickly pushing an application to the repository.`,
	Args:  cobra.MaximumNArgs(3),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("must provide three arguments <REPOSITORY ORG/BUCKET/APP SOURCE>")
		}
		words := strings.Split(args[1], "/")
		if len(words) < 3 {
			return fmt.Errorf("invalid format for <org/bucket/app> argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		urlPath := args[0]
		repoPath := strings.Split(args[1], "/")
		buildablePath := args[2]

		pkgBuilder, err := getPackageBuilder("BUILDABLE", buildablePath)
		if err != nil {
			SetError(err, 2)
			return
		}
		defer pkgBuilder.Close()

		err = modifyPackageBuilder(pkgBuilder)
		if err != nil {
			SetError(err, 3)
			return
		}

		err = pushPackage(pkgBuilder, urlPath, repoPath)
		if err != nil {
			SetError(err, 4)
			return
		}

		return
	},
}

func init() {
	f := pushCmd.Flags()
	f.StringVarP(&flagKey, "key", "k", "", "vrepo authentication key file name")
}

// checkAuthentication checks to see if the flag has provided a name if not
// checks to see if default exists if that doesnt exist
// errors out saying you need to provide authentication
func checkAuthentication() (string, error) {
	usr, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	pathCheck := filepath.Join(usr, ".vorteil", "repository-keys")
	token, err := checkDefaultAndProvided(pathCheck)
	if err != nil {
		return "", err
	}

	return token, nil
}

// checkDefaultAndProvided checks both files for authentication
func checkDefaultAndProvided(pathCheck string) (string, error) {
	var token string
	var path string
	var err error
	if flagKey != "" {
		path = filepath.Join(pathCheck, flagKey)
		token, err = checkAuthFile(path)
		if err != nil {
			return "", fmt.Errorf("unable to locate '%s' keyfile", flagKey)
		}
	} else {
		path = filepath.Join(pathCheck, "default")
		token, err = checkAuthFile(path)
		if err != nil {
			return "", errors.New("unable to find any authentication provided try using the --key flag")
		}
	}
	return token, nil
}

// checkAuthFile checks if file exists returns token or error
func checkAuthFile(path string) (string, error) {
	_, err := os.Stat(path)
	if err == nil {
		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()
		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
	return "", fmt.Errorf("unable to find keyfile at %s", path)
}

// preparePackage for upload
func preparePackage(builder vpkg.Builder) (*os.File, error) {
	spinner := log.NewProgress("Preparing Package", "", 0)
	defer spinner.Finish(true)
	file, err := ioutil.TempFile(os.TempDir(), "vpkg-")
	if err != nil {
		return nil, err
	}

	err = builder.Pack(file)
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// generateRequest creates request to send to the repository
func generateRequest(url string, repo []string, r io.ReadCloser, token string) (*http.Request, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/organisations/%s/buckets/%s/apps/%s", url, repo[0], repo[1], repo[2]), r)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	return req, nil
}

// uploadPackage sends the request to upload package
func uploadPackage(url string, repo []string, token string, file *os.File) error {
	client := &http.Client{}

	stats, err := file.Stat()
	if err != nil {
		return err
	}

	p := log.NewProgress("Uploading Package", "KiB", stats.Size())
	r := p.ProxyReader(file)
	defer p.Finish(true)

	req, err := generateRequest(url, repo, r, token)
	if err != nil {
		p.Finish(false)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		p.Finish(false)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.Finish(false)
			return err
		}
		return fmt.Errorf("%s: %s", errors.New(resp.Status), string(bodyBytes))
	}
	return nil
}

// pushPackage takes builder, url and repo array of strings which is org/bucket/app
func pushPackage(builder vpkg.Builder, url string, repo []string) error {

	// check authentication before doing things
	token, err := checkAuthentication()
	if err != nil {
		return err
	}

	file, err := preparePackage(builder)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	err = uploadPackage(url, repo, token, file)
	if err != nil {
		return err
	}

	return nil
}