package folder

import (
	"fmt"
	"github.com/gofrs/uuid"
)

/* Error checkings
	1. Moving a node to a child of itself DONE
	2. Moving a folder to itself DONE
	3. Moving folders across orgId DONE
	4. Invalid source folder name DONE
	5. Invalid dest folder name DONE
*/

// Assumption: 
// There can only be one orgId that contains 
// a specific pair of folders with the names "name1" and "name2". 
func (d *driver) MoveFolder(name string, dst string) ([]Folder, error) {

	if name == dst {
		return nil, fmt.Errorf("moveFolder: cannot move folder to itself")
	}
	// validate folder names

	// check for existence, obtain their folder pointers
	srcNameFolders := d.nameIndex[name]
	destNameFolders := d.nameIndex[dst]

	// temp map to store all the seen orgIds with sourceName
	srcNameOrgIds := make(map[uuid.UUID]bool)
	for _, folder := range srcNameFolders {
		srcNameOrgIds[folder.OrgId] = true
	}
	var matchingOrgId uuid.UUID
	orgIdMatched := false
	// look for a folder with dest name, and orgId
	for _, folder := range destNameFolders {
		// check if there exists a srcName folder with orgId of the destNameFolder
		_, found := srcNameOrgIds[folder.OrgId]
		// if found, then its a valid match with folder.OrgId, destName and srcName
		if found {
			matchingOrgId = folder.OrgId
			orgIdMatched = true
		}
	}
	if !orgIdMatched {
		return nil, fmt.Errorf("moveFolder: No matching orgId found between '%s' and '%s'", name, dst)
	}

	// obtain folder pointers for name and dst
	var srcFolder *Folder
	var dstFolder *Folder
	for _, folder := range d.orgIdIndex[matchingOrgId] {
		f := folder
		if f.Name == name {
			srcFolder = f
		} else if f.Name == dst {
			dstFolder = f
		}
	}
	if srcFolder == nil {
		return nil, fmt.Errorf("moveFolder: Source folder '%s' not found", name)
	} else if dstFolder == nil {
		return nil, fmt.Errorf("moveFolder: Destination folder '%s' not found", dst)
	}

	// check for cycles, dst cannot be a descendent of src
	if (d.isDescendent(dstFolder, srcFolder)) {
		return nil, fmt.Errorf("moveFolder: dstFolder '%s' cannot be a descendent of scrFolder '%s'", dst, name)
	}

	// if srcFolder is already in the dstFolder, do nothing
	if srcFolder.Parent == dstFolder {
		return d.getAllFolders(), nil
	}

	// remove srcFolder from its current parent
	if srcFolder.Parent != nil {
		err := d.removeChild(srcFolder.Parent, srcFolder)
		if err != nil {
			return nil, fmt.Errorf("moveFolder: error removing source folder from its current parent: %v", err)
		}
	}

	dstFolder.Children = append(dstFolder.Children, srcFolder)
	srcFolder.Parent = dstFolder

	// recursive reassign children's path
	// oldPath := srcFolder.Paths
	newPath := fmt.Sprintf("%s.%s", dstFolder.Paths, name)
	fmt.Printf("New path for folder %s is %s\n", srcFolder.Paths, newPath)
	
	// update all the paths and update pathIndex along the way
	d.updatePaths(srcFolder, newPath)

	return d.getAllFolders(), nil
}

// recursively updates all path names
func (d *driver) updatePaths(currFolder *Folder, newPath string) {
	oldPaths := currFolder.Paths
	currFolder.Paths = newPath
	fmt.Printf("old currFolder.Paths = %s, new = %s\n", oldPaths, newPath)
	delete(d.pathIndex, oldPaths)

	// remove oldPath index, create new pathIndex
	d.pathIndex[newPath] = currFolder		

	for _, child := range currFolder.Children {
		c := child
		newPath := fmt.Sprintf("%s.%s", newPath, c.Name)
		d.updatePaths(c, newPath)
	}
} 

// Checks if child is the descendent of the parent
// recrusively bubble up child folder
func (d *driver) isDescendent(child, parent *Folder) bool {
	if child == parent {
		return true
	}
	if child.Parent == nil {
		return false
	}
	return d.isDescendent(child.Parent, parent)
}

// removes a child folder from a parent's children slice.
func (d *driver) removeChild(parent *Folder, child *Folder) error {
	for i, c := range parent.Children {
		if c == child {
			// Remove child from slice
			parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("removeChild: child folder '%s' not found under parent '%s'", child.Name, parent.Name)
}


