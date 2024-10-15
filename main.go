package main

import (
	"fmt"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

func main() {
	// orgID := uuid.FromStringOrNil(folder.DefaultOrgID)
	orgId1 := uuid.Must(uuid.NewV4())
	// orgId2 := uuid.Must(uuid.NewV4())
	// folders := folder.GetAllFolders()
	var folders = []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgId1},
				{Name: "delta", Paths: "delta", OrgId: orgId1},
	}

	// example usage
	folderDriver, err := folder.NewDriver(folders)

	if err != nil {
		fmt.Printf("Error when creating driver: %v", err)
		return
	}

	f, err := folderDriver.MoveFolder("bravo", "delta")
	if err != nil {
		fmt.Printf("Error when creating driver: %v", err)
	}

	// folder.PrettyPrint(folders)
	// fmt.Printf("\n Folders for orgID: %s", orgID)
	folder.PrettyPrint(f)
}
