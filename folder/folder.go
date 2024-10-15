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

	// Component 1
	GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error)

	// Component 2
	MoveFolder(name string, dst string) ([]Folder, error)
}

// Manages folder hierarchy
type driver struct {
	// Stores slice of all folders
	folders   []*Folder
	pathIndex map[string]*Folder
	// nameIndex and orgId index must account for folders with duplcate names and orgIds
	nameIndex  map[string][]*Folder
	orgIdIndex map[uuid.UUID][]*Folder
	mu         sync.RWMutex
}

// Initialises FolderDriver, populating parent child
// TODO: Implement cycle detection
func NewDriver(folders []Folder) (IDriver, error) {
	folderDriver := &driver{
		folders:    []*Folder{},
		pathIndex:  make(map[string]*Folder),
		nameIndex:  make(map[string][]*Folder),
		orgIdIndex: make(map[uuid.UUID][]*Folder),
		// mutex lock does not require explicit initialisation
	}

	// populate folders and maps of driver, folders is a slice of Folder, NOT *Folder
	for _, folder := range folders {
		f := folder

		// Check for duplicate folder names within the same OrgId
		for _, existingFolder := range folderDriver.nameIndex[f.Name] {
			if existingFolder.OrgId == f.OrgId {
				return nil, fmt.Errorf("newDriver: duplicate folder name '%s' in OrgId '%s'", f.Name, f.OrgId)
			}
		}

		// Check for cycles
		if hasRepeats(f.Paths) {
			return nil, fmt.Errorf("newDriver: cannot instantiate path %s has it will create a cycle", folder.Paths)
		}

		folderDriver.folders = append(folderDriver.folders, &f)
		folderDriver.pathIndex[f.Paths] = &f
		folderDriver.nameIndex[f.Name] = append(folderDriver.nameIndex[f.Name], &f)
		folderDriver.orgIdIndex[f.OrgId] = append(folderDriver.orgIdIndex[f.OrgId], &f)
	}

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
			return nil, fmt.Errorf("newDriver: Parent oath '%s' not found for folder '%s'", parentPath, folder.Name)
		}

		// establish parent child
		parentFolder.Children = append(parentFolder.Children, folder)
		folder.Parent = parentFolder
	}

	return folderDriver, nil
}

// get the substring of childPath up until the last dot
func getParentPath(childPath string) string {
	lastDot := strings.LastIndex(childPath, ".")
	if lastDot == -1 {
		return ""
	}
	return childPath[:lastDot]
}

// Checks if there exists repeated foldeName in path string
func hasRepeats(s string) bool {
	parts := strings.Split(s, ".")
	seen := make(map[string]bool)
	for _, part := range parts {
		if seen[part] {
			return true
		}
		seen[part] = true
	}
	return false
}

// Returns referenced folders
func (d *driver) getAllFolders() []Folder {
	var allFolders []Folder
	for _, folder := range d.folders {
		allFolders = append(allFolders, *folder)
	}
	return allFolders
}
