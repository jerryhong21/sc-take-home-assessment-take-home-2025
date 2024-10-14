package folder_test

import (
	"reflect"
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)


// func Test_folder_GetFoldersByOrgID(t *testing.T) {
// 	t.Parallel()
// 	tests := [...]struct {
// 		name    string
// 		orgID   uuid.UUID
// 		folders []folder.Folder
// 		want    []folder.Folder
// 	}{

// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			f, error := folder.NewDriver(tt.folders)
// 			if error == nil {
// 				t.Errorf("Creating drier")
// 				return
// 			}
// 			get := f.GetFoldersByOrgID(tt.orgID)


// 		})
// 	}
// }

func TestGetAllChildFolders(t *testing.T) {
	orgID1 := uuid.Must(uuid.NewV4())
	orgID2 := uuid.Must(uuid.NewV4())

	tests := []struct {
		name         string
		orgID        uuid.UUID
		parent       string
		want         []folder.Folder
		wantErr      bool
		// extra folders for specific scenarios like cyclic paths or duplicates to add at each test
		extraFolders []folder.Folder
	}{

		// basic function tests
		{
			name:   "Retrieving child folders from a parent folder that has no children",
			orgID:  orgID1,
			parent: "echo",
			want:   []folder.Folder{},
			wantErr: false,
		},
		{
			name:  "parent with multiple direct and indirect child folders",
			orgID: orgID1,
			parent: "alpha",
			want: []folder.Folder{
				{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
				{Name: "charlie", OrgId: orgID1, Paths: "alpha.bravo.charlie"},
				{Name: "delta", OrgId: orgID1, Paths: "alpha.delta"},
			},
			wantErr: false,
		},
		{
			name:   "root nodes without children return an empty list",
			orgID:  orgID1,
			parent: "echo",
			want:   []folder.Folder{},
			wantErr: false,
		},
		{
			name:  "nested child folders with large depth",
			orgID: orgID1,
			parent: "bravo",
			want: []folder.Folder{
				{Name: "charlie", OrgId: orgID1, Paths: "alpha.bravo.charlie"},
			},
			wantErr: false,
		},
		{
			name:  "multiple parent folders sharing the same name within the same org",
			orgID: orgID1,
			parent: "commonParent",
			want: nil,
			wantErr: true,
			extraFolders: []folder.Folder{
				// Create two distinct parents named commonParent with unique paths.
				{Name: "commonParent", OrgId: orgID1, Paths: "commonParent1"},
				{Name: "commonParent", OrgId: orgID1, Paths: "commonParent2"},
				{Name: "child1", OrgId: orgID1, Paths: "commonParent1.child1"},
				{Name: "child2", OrgId: orgID1, Paths: "commonParent2.child2"},
			},
		},
		{
			name:    "duplicate folder names under the same parent and org",
			orgID:   orgID1,
			parent:  "alpha",
			want:    nil, // Expecting an error due to duplicate names.
			wantErr: true,
			extraFolders: []folder.Folder{
				// duplicate folder named "bravo" under "alpha"
				{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
			},
		},
		{
			name:    "case sensitive parent names",
			orgID:   orgID1,
			parent:  "Alpha",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "Managing multiple folders with identical names across different organizations",
			orgID: orgID2,
			parent: "foxtrot",
			want: []folder.Folder{
				{Name: "golf", OrgId: orgID2, Paths: "foxtrot.golf"},
			},
			wantErr: false,
		},

		// edge cases
		{
			name:    "Attempting to retrieve child folders from a non-existent parent folder",
			orgID:   orgID1,
			parent:  "invalid_folder",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Trying to access a parent folder that belongs to a different org",
			orgID:   orgID1,
			// 'foxtrot' already exists under OrgID2, not 1
			parent:  "foxtrot", 
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Handling scenarios where the provided OrgID is invalid",
			orgID:   uuid.Nil, 
			parent:  "alpha",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Ensuring proper error handling when the parent name is empty",
			orgID:   orgID1,
			parent:  "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Handling parent names containing Unicode characters gracefully",
			orgID:   orgID1,
			// non ascii characters
			parent:  "αlphα",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "parent folder is a leaf node with no children",
			orgID:   orgID1,
			parent:  "echo",
			want:    []folder.Folder{},
			wantErr: false,
		},
		{
			name:    "invalid folder paths strings",
			orgID:   orgID1,
			parent:  "alpha..bravo",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "parent names which are substrings of existing folders should not get matches (false-positives)",
			orgID:   orgID1,
			parent:  "alp",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "detect cycles within folder paths to avoid infinite loops",
			orgID:   orgID1,
			parent:  "alpha.charlie",
			want:    nil,
			wantErr: true,
			extraFolders: []folder.Folder{
				// cycle folder path : alpha.charlie.alpha -> cycles back to alpha.
				{Name: "alpha.charlie.alpha", OrgId: orgID1, Paths: "alpha.charlie.alpha"},
			},
		},
	}

	for _, tt := range tests {
		// if we don't do this, parallel tests may mess up pointer
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			
			folders := []folder.Folder{
				{Name: "alpha", OrgId: orgID1, Paths: "alpha"},
				{Name: "bravo", OrgId: orgID1, Paths: "alpha.bravo"},
				{Name: "charlie", OrgId: orgID1, Paths: "alpha.bravo.charlie"},
				{Name: "delta", OrgId: orgID1, Paths: "alpha.delta"},
				{Name: "echo", OrgId: orgID1, Paths: "echo"},

				// Folders belonging to OrgID2.
				{Name: "foxtrot", OrgId: orgID2, Paths: "foxtrot"},
				{Name: "golf", OrgId: orgID2, Paths: "foxtrot.golf"},
			}

			// add additional folders required for the specific test case.
			if len(tt.extraFolders) > 0 {
				folders = append(folders, tt.extraFolders...)
			}
			driver, err := folder.NewDriver(folders)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewDriver() error = %v, expected error presence: %v", err, tt.wantErr)
				return
			}

			// If an error was expected during driver initialization, no further testing is needed.
			if tt.wantErr {
				return
			}

			got, err := driver.GetAllChildFolders(tt.orgID, tt.parent)

			// error checking
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllChildFolders() error = %v, expected error presence: %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Errorf("GetAllChildFolders() expected an error but none was returned")
				return
			}

			// deep equal
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllChildFolders() returned %v, but expected %v", got, tt.want)
			}
		})
	}
}