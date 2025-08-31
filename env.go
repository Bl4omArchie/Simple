package simple


import (
	"os"
	
	"github.com/joho/godotenv"
)


// Open a .env file
// Return results in a slice of string
func OpenEnv(tags ...string) []string {
	err := godotenv.Load()
	if err != nil {
		return nil
	}

	var tagsSlice []string
	for _, tag := range tags {
		tagsSlice = append(tagsSlice, os.Getenv(tag))
	}
	
	return tagsSlice
}

// Open a .env file with one or more filename
// Return results in a slice of string
func OpenEnvFilenames(filenames []string, tags ...string) []string {
	err := godotenv.Load(filenames ...)
	if err != nil {
		return nil
	}

	var tagsSlice []string
	for _, tag := range tags {
		tagsSlice = append(tagsSlice, os.Getenv(tag))
	}
	
	return tagsSlice
}
