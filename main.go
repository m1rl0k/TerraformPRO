package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"path/filepath"
)

func main() {

    // Check if the "-cleanup" argument was provided
    if len(os.Args) > 1 && os.Args[1] == "-cleanup" {
        // Call the cleanup function
        err := cleanup()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // Print a message to indicate that cleanup is complete
        fmt.Println("Cleanup complete.")
        return
    }

    // Prompt the user to select a provider
    fmt.Print("Which provider do you want to use? (azure/aws/google) ")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Scan()
    provider := strings.ToLower(scanner.Text())

    // Check if the selected provider is valid
    switch provider {
    case "azure", "aws", "google":
        // Navigate to the appropriate directory
        dir := fmt.Sprintf("./%s/templates", provider)
        err := os.Chdir(dir)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        // Print a list of available templates in the templates directory
        fmt.Printf("Available templates for %s:\n", provider)
        templates, err := listDirectories(".")
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        for i, template := range templates {
            fmt.Printf("%d. %s\n", i+1, template)
        }

        // Prompt the user to select a template
        fmt.Print("Which template do you want to use? ")
        scanner.Scan()
        templateNum := scanner.Text()

        // Check if the selected template is valid
        if templateNumInt, err := strconv.Atoi(templateNum); err == nil && templateNumInt > 0 && templateNumInt <= len(templates) {
            selectedTemplate := templates[templateNumInt-1]

            // Navigate to the selected template's directory
            err := os.Chdir(selectedTemplate)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            // If the provider is "azure", prompt the user to enter a subscription ID
            if provider == "azure" {
                fmt.Print("Enter your subscription ID: ")
                scanner.Scan()
                subscriptionID := scanner.Text()

                // Update the subscription ID in the provider.tf file
                err = updateSubscriptionID(subscriptionID)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }

                fmt.Printf("Updated subscription ID to %s in provider.tf\n", subscriptionID)

                // Wait for 5 seconds after updating the subscription ID
                fmt.Println("Waiting for 5 seconds...")
                time.Sleep(5 * time.Second)
            } else if provider == "aws" {
                // Prompt the user to enter the AWS region, access key, and secret key
                fmt.Print("Enter your AWS region: ")
                scanner.Scan()
                region := scanner.Text()

                fmt.Print("Enter your AWS access key ID: ")
                scanner.Scan()
                accessKey := scanner.Text()

                fmt.Print("Enter your AWS secret access key: ")
                scanner.Scan()
                secretKey := scanner.Text()

                // Update the AWS provider configuration in the provider.tf file
                err = updateAWSProviderConfig(region, accessKey, secretKey)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }

                fmt.Printf("Updated AWS provider configuration in provider.tf\n")

                // Wait for 5 seconds after updating the AWS provider configuration
                fmt.Println("Waiting for 5 seconds...")
                time.Sleep(5 * time.Second)
				} else if provider == "google" {
               // Prompt the user to enter the Google Cloud project ID and region
			   fmt.Print("Enter your Google Cloud project ID: ")
               scanner.Scan()
               project := scanner.Text()

			   fmt.Print("Enter your Google Cloud region: ")
               scanner.Scan()
               region := scanner.Text()

               // Update the Google Cloud provider configuration in the provider.tf file
			   err = updateGoogleProviderConfig(project, region)
               if err != nil {
               fmt.Println(err)
               os.Exit(1)

}

    fmt.Printf("Updated Google Cloud provider configuration in provider.tf\n")

    // Wait for 5 seconds after updating the Google Cloud provider configuration
    fmt.Println("Waiting for 5 seconds...")
    time.Sleep(5 * time.Second)
}
			
			// Run "terraform init" in the selected template's directory
			cmd := exec.Command("terraform", "init")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Initialized Terraform provider for %s in template %s\n", provider, selectedTemplate)
		} else {
			fmt.Println("Invalid template selected.")
			os.Exit(1)
		}
	default:
		fmt.Println("Invalid provider selected.")
		os.Exit(1)
	}
}
// listDirectories returns a slice of the names of all directories in the
// specified path.
func listDirectories(path string) ([]string, error) {
	var dirs []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}
	return dirs, nil
}
// updateSubscriptionID updates the subscription ID in the providers.tf file
// in the current directory.
func updateSubscriptionID(subscriptionID string) error {
	// Read the contents of the providers.tf file
	providerFile, err := ioutil.ReadFile("provider.tf")
	if err != nil {
		return err
	}

	// Find the subscription ID line in the file and replace its value with the new subscription ID
	newProviderFile := ""
	lines := strings.Split(string(providerFile), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "  subscription_id") {
			if strings.Contains(line, "\"") {
				line = fmt.Sprintf("  subscription_id = \"%s\"", subscriptionID)
			} else {
				line = fmt.Sprintf("  subscription_id = %s", subscriptionID)
			}
		}
		newProviderFile += line + "\n"
	}

	// Write the updated providers.tf file back to disk
	err = ioutil.WriteFile("provider.tf", []byte(newProviderFile), 0644)
	if err != nil {
		return err
	}

	return nil
}
func updateAWSProviderConfig(region string, accessKey string, secretKey string) error {
	// Read the contents of the providers.tf file
	providerFile, err := ioutil.ReadFile("provider.tf")
	if err != nil {
		return err
	}

	// Update the AWS provider configuration in the file with the new values
	newProviderFile := ""
	inProviderBlock := false
	lines := strings.Split(string(providerFile), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "provider \"aws\"") {
			inProviderBlock = true
		}
		if inProviderBlock {
			if strings.Contains(line, "region") {
				line = fmt.Sprintf("  region = \"%s\"", region)
			}
			if strings.Contains(line, "access_key") {
				line = fmt.Sprintf("  access_key = \"%s\"", accessKey)
			}
			if strings.Contains(line, "secret_key") {
				line = fmt.Sprintf("  secret_key = \"%s\"", secretKey)
			}
		}
		if strings.HasPrefix(line, "}") {
			inProviderBlock = false
		}
		newProviderFile += line + "\n"
	}

	// Write the updated providers.tf file back to disk
	err = ioutil.WriteFile("provider.tf", []byte(newProviderFile), 0644)
	if err != nil {
		return err
	}

	return nil
}
// updateGoogleProviderConfig updates the Google Cloud provider configuration in the providers.tf file
// in the current directory.
func updateGoogleProviderConfig(project string, region string) error {
    // Read the contents of the providers.tf file
    providerFile, err := ioutil.ReadFile("provider.tf")
    if err != nil {
        return err
    }

    // Update the Google Cloud provider configuration in the file with the new values
    newProviderFile := ""
    inProviderBlock := false
    lines := strings.Split(string(providerFile), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "provider \"google\"") {
            inProviderBlock = true
        }
        if inProviderBlock {
            if strings.Contains(line, "project") {
                line = fmt.Sprintf("  project = \"%s\"", project)
            }
            if strings.Contains(line, "region") {
                line = fmt.Sprintf("  region = \"%s\"", region)
            }
        }
        if strings.HasPrefix(line, "}") {
            inProviderBlock = false
        }
        newProviderFile += line + "\n"
    }

    // Write the updated providers.tf file back to disk
    err = ioutil.WriteFile("provider.tf", []byte(newProviderFile), 0644)
    if err != nil {
        return err
    }

    return nil
}
func cleanup() error {
    // Find all .terraform directories and .terraform.lock.hcl files recursively
    terraformDirs := []string{}
    lockFiles := []string{}
    err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() && info.Name() == ".terraform" {
            terraformDirs = append(terraformDirs, path)
        } else if !info.IsDir() && info.Name() == ".terraform.lock.hcl" {
            lockFiles = append(lockFiles, path)
        }
        return nil
    })
    if err != nil {
        return err
    }

    // Remove all .terraform directories and .terraform.lock.hcl files
    for _, dir := range terraformDirs {
        err = os.RemoveAll(dir)
        if err != nil {
            return err
        }
    }
    for _, file := range lockFiles {
        err = os.Remove(file)
        if err != nil {
            return err
        }
    }

    return nil
}