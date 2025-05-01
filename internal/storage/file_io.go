package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileIO interface {
	CreateRepoInitInfo(repositoryId int) error
	GetRepositoryId() (int, error)
}

type FileIOImpl struct {
	currentDirectory string
	infoFile         string
	infoDir          string
}

func NewFileIO() *FileIOImpl {
	currDirr, err := os.Getwd()
	if err != nil {
		panic("please give oma execution rights")
	}

	return &FileIOImpl{
		currentDirectory: currDirr,
		infoFile:         "repository_info.txt",
		infoDir:          ".oma",
	}
}

func (repoFileIO *FileIOImpl) CreateRepoInitInfo(repositoryId int) error {
	// permissions in Go are in octal notation apparently
	// hence the 0 prefix
	err := os.Mkdir(".oma", 0755)
	if err != nil {
		return fmt.Errorf("error while creating repository info file parent:\n%v", err)
	}

	file, err := os.Create(filepath.Join(repoFileIO.currentDirectory, repoFileIO.infoDir, repoFileIO.infoFile))
	defer file.Close()

	if err != nil {
		return fmt.Errorf("error while creating repository info file:\n%v", err)
	}

	repoIdInBytes := []byte("repositoryId=" + strconv.Itoa(repositoryId))
	file.Write(repoIdInBytes)

	return nil
}

func (repoFileIO *FileIOImpl) GetRepositoryId() (int, error) {
	contentBytes, err := os.ReadFile(filepath.Join(repoFileIO.infoDir, repoFileIO.infoFile))

	if err != nil {
		return -1, fmt.Errorf("error while reading the info file:\n%v", err)
	}

	content := strings.ReplaceAll(string(contentBytes), "\r", "")
	content = strings.ReplaceAll(content, "\t", "  ")

	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "repositoryId") {
			repoIdString := strings.Split(line, "=")[1]
			repoId, err := strconv.ParseInt(repoIdString, 10, 64)

			if err != nil {
				return -1, fmt.Errorf("error while parsing repository ID:\n%v", err)
			}

			return int(repoId), nil
		}
	}

	return -1, fmt.Errorf("repository ID could not be found")
}
