package folder

import (
	"fmt"
	"github.com/gofrs/uuid"
	
)


func GetAllFolders() []Folder {
	return GetSampleData()
}

func (d *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	folders := d.folders

	res := []Folder{}
	for _, f := range folders {
		if f.OrgId == orgID {
			res = append(res, *f)
		}
	}

	return res
}

func (driver *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
	// acquire and release read lock
	driver.mu.RLock()
	defer driver.mu.RUnlock()

	if orgID == uuid.Nil {
		return nil, fmt.Errorf("GetAllChildFolders: invalid OrgID - '%s'", orgID)
	}

	// get a list of names from the map and filter by orgId
	folders, found := driver.nameIndex[name]
	if !found {
		return nil, fmt.Errorf("GetAllChildFolders: folder '%s' does not exist", name)
	}

	// check if duplicate names exist in this
	orgIdFound := false
	var parentFolder *Folder
	for _, folder := range folders {
		if folder.OrgId == orgID {
			if orgIdFound {
				return nil, fmt.Errorf("GetAllChildrenFolders: there exists more than one '%s' in org %s", name, orgID)
			}
			orgIdFound = true
			parentFolder = folder
		}
	}

	if !orgIdFound {
		return nil, fmt.Errorf("GetAllChildFolders: folder '%s' does not exist in org '%s'", name, orgID)
	}

	allChildren, err := collectAllDescendents(parentFolder)
	if err != nil {
		return nil, fmt.Errorf("GetAllChildrenFolders: there exists more than one '%s' in org '%s'", name, orgID)
	}
	if len(allChildren) == 0 {
		fmt.Printf("GetAllChildrenFolders: '%s' has no children folders\n", name)
	}

	return allChildren, nil
}

// Recursive function to collect all children
func collectAllDescendents(parent *Folder) ([]Folder, error) {

	allChildren := []Folder{}
	var recursiveCollect func(*Folder)
	recursiveCollect = func(currFolder *Folder) {
		for _, child := range currFolder.Children {
			allChildren = append(allChildren, *child)
			recursiveCollect(child)
		}
	}

	recursiveCollect(parent)
	return allChildren, nil
}
