package main

import (
	"fmt"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)

func main() {
	// orgID := uuid.FromStringOrNil(folder.DefaultOrgID)
	orgId1 := uuid.Must(uuid.NewV4())
	orgId2 := uuid.Must(uuid.NewV4())
	// folders := folder.GetAllFolders()
	var folders = []folder.Folder{
		{Name: "alpha", OrgId: orgId1, Paths: "alpha"},
		{Name: "bravo", OrgId: orgId1, Paths: "alpha.bravo"},
		{Name: "charlie", OrgId: orgId1, Paths: "alpha.bravo.charlie"},
		{Name: "delta", OrgId: orgId1, Paths: "alpha.delta"},
		{Name: "echo", OrgId: orgId1, Paths: "echo"},

		// Folders belonging to OrgID2.
		{Name: "foxtrot", OrgId: orgId2, Paths: "foxtrot"},
		{Name: "golf", OrgId: orgId2, Paths: "foxtrot.golf"},
		{Name: "commonParent", OrgId: orgId1, Paths: "commonParent1"},
		{Name: "commonParent", OrgId: orgId1, Paths: "commonParent2"},
		{Name: "child1", OrgId: orgId1, Paths: "commonParent1.child1"},
		{Name: "child2", OrgId: orgId1, Paths: "commonParent2.child2"},
	}

	// example usage
	folderDriver, err := folder.NewDriver(folders)
	if err != nil {
		fmt.Printf("Error when creating driver: %v", err)
		return
	}

	childFolder, err := folderDriver.GetAllChildFolders(orgId1, "commonParent")
	if err != nil {
		fmt.Printf("Error when creating driver: %v", err)
	}

	// folder.PrettyPrint(folders)
	// fmt.Printf("\n Folders for orgID: %s", orgID)
	folder.PrettyPrint(childFolder)
}
