package folder

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
)

type IDriver interface {
	// GetFoldersByOrgID returns all folders that belong to a specific orgID.
	GetFoldersByOrgID(orgID uuid.UUID) []Folder
	// component 1
	// Implement the following methods:
	// GetAllChildFolders returns all child folders of a specific folder.
	GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error)

	// component 2
	// Implement the following methods:
	// MoveFolder moves a folder to a new destination.
	MoveFolder(name string, dst string) ([]Folder, error)
}

// Manages folder hierarchy
type driver struct {
	folders     []*Folder // Stores slice of all folders
	pathIndex   map[string]*Folder
	// nameIndex and orgId index must account for folders with duplcate names and orgIds
	nameIndex   map[string][]*Folder
	orgIdIndex  map[uuid.UUID][]*Folder
	mu          sync.RWMutex
}

// Initialises FolderDriver, populating parent child
func NewDriver(folders []Folder) (IDriver, error) {
	folderDriver := &driver{
		folders: []*Folder{},
		pathIndex: make(map[string]*Folder),
		nameIndex: make(map[string][]*Folder),
		orgIdIndex: make(map[uuid.UUID][]*Folder),
		// mutex lock does not require explicit initialisation
	}

	// populate folders and maps of driver
	// here, folders is a slice of Folder, NOT *Folder
	for _, folder := range folders {
		f := folder

		// Check for duplicate folder names within the same OrgId
		for _, existingFolder := range folderDriver.nameIndex[f.Name] {
			if existingFolder.OrgId == f.OrgId {
				return nil, fmt.Errorf("duplicate folder name '%s' in OrgId '%s'", f.Name, f.OrgId)
			}
		}

		folderDriver.folders = append(folderDriver.folders, &f)
		folderDriver.pathIndex[f.Paths] = &f
		folderDriver.nameIndex[f.Name] = append(folderDriver.nameIndex[f.Name], &f)
		folderDriver.orgIdIndex[f.OrgId] = append(folderDriver.orgIdIndex[f.OrgId], &f)
		fmt.Printf("Processed %s\n", folder.Name)
	}

	// precondition: folderDriver.folder and all folderDriver maps are populated
	// folderDriver.folders stores a slice of *Folder
	for _, folder := range folderDriver.folders {
		// populate child and parent
		// obtain parent and set children
		parentPath := getParentPath(folder.Paths)
		// root path
		if parentPath == "" {
			continue
		}
		parentFolder, found := folderDriver.pathIndex[parentPath]
		if !found {
			fmt.Printf("Error: Parent oath '%s' not found for folder '%s'\n", parentPath, folder.Name)
			continue
		}
	
		// establish parent child
		parentFolder.Children = append(parentFolder.Children, folder)
		folder.Parent = parentFolder
	}

	return folderDriver, nil
}

func getParentPath(childPath string) string {
	lastDot := strings.LastIndex(childPath, ".")
	if lastDot == -1 {
		return ""
	}
	return childPath[:lastDot]
}

// type driver struct {
// 	// define attributes here
// 	// data structure to store folders
// 	// or preprocessed data

// 	*FolderDriver
// }

// // receives folders, and returns a folderDriver
// func NewDriver(folders []Folder) IDriver {
// 	folderDriver := InitialiseFolderDriver(folders)
// 	return &driver{folderDriver}
// }
