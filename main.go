package main

import (
	"fmt"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

func main() {
	orgID := uuid.FromStringOrNil(folder.DefaultOrgID)

	folders := folder.GetAllFolders()

	// example usage
	folderDriver := folder.NewDriver(folders)
	orgFolder := folderDriver.GetFoldersByOrgID(orgID)

	folder.PrettyPrint(folders)
	fmt.Printf("\n Folders for orgID: %s", orgID)
	folder.PrettyPrint(orgFolder)
}
