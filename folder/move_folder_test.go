package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
)


func Test_folder_MoveFolder(t *testing.T) {

	orgId1 := uuid.Must(uuid.NewV4())
	orgId2 := uuid.Must(uuid.NewV4())

	tests := []struct {
		name          string
		initialFolders []folder.Folder
		source        string
		destination   string
		want          []folder.Folder
		wantRunError   bool
	}{
		// Functionalities
		{
			name: "1. Moving folder with no children into another folder with no children",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: orgId2},
			},
			source:      "bravo",
			destination: "alpha",
			// same thing
			want: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: orgId2},
			},
			wantRunError: false,
		},

		{
			name: "2. moving folder with no children into another folder with children",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgId1},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgId1},
				{Name: "bravo", Paths: "bravo", OrgId: orgId1},
			},
			source:      "bravo",
			destination: "delta",
			want: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgId1},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgId1},
				// move into detlta
				{Name: "bravo", Paths: "alpha.delta.bravo", OrgId: orgId1},
			},
			wantRunError: false,
		},
		{
			name: "3. Moving folder with children into another folder with no children",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgId1},
				{Name: "delta", Paths: "delta", OrgId: orgId1},
			},
			source:      "bravo",
			destination: "delta",
			want: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "delta", Paths: "delta", OrgId: orgId1},
				{Name: "bravo", Paths: "delta.bravo", OrgId: orgId1},
				{Name: "charlie", Paths: "delta.bravo.charlie", OrgId: orgId1},
			},
			wantRunError: false,
		},
		{
			name: "4. Moving folder with children into another folder with children",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				// with children
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgId1},
				// with children
				{Name: "delta", Paths: "alpha.delta", OrgId: orgId1},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgId1},
			},
			source: "bravo",
			destination: "delta",
			want: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				// delta becomes 4 nested
				{Name: "delta", Paths: "alpha.delta", OrgId: orgId1},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.delta.bravo", OrgId: orgId1},
				{Name: "charlie", Paths: "alpha.delta.bravo.charlie", OrgId: orgId1},
			},
			wantRunError: false,
		},
		{
			name: "5. Move a deeply nested folder into another folder with children",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgId1},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgId1},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgId1},
				{Name: "foxtrot", Paths: "alpha.bravo.charlie.foxtrot", OrgId: orgId1},
			},
			source: "charlie",
			destination: "delta",
			want: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgId1},
				{Name: "echo", Paths: "alpha.delta.echo", OrgId: orgId1},
				{Name: "charlie", Paths: "alpha.delta.charlie", OrgId: orgId1},
				{Name: "foxtrot", Paths: "alpha.delta.charlie.foxtrot", OrgId: orgId1},
			},
			wantRunError: false,
		},
		{
			name: "6. Move folder in-place (move to its current parent)",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
			},
			source: "bravo",
			destination: "alpha",
			want: []folder.Folder{
				// no change expected
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
			},
			wantRunError: false,
		},
		{
			name: "7. Move folder when source is root folder",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "bravo", OrgId: orgId1},
				{Name: "charlie", Paths: "charlie", OrgId: orgId1},
			},
			source:      "bravo",
			destination: "alpha",
			want: []folder.Folder{
				// bravo should move into alpha, charlie should stay the same
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
				{Name: "charlie", Paths: "charlie", OrgId: orgId1},
			},
			wantRunError: false,
		},

		// edge cases an error checking

		{
			name: "8. error moving a node to a child of itself",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: orgId1},
			},
			source: "bravo",
			destination: "charlie",
			want: nil,
			wantRunError: true,
		},
		{
			name: "9. error when moving a folder to itself",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
			},
			source:         "bravo",
			destination:     "bravo",
			want:            nil,
			wantRunError:    true,
		},
		{
			name: "10. error: moving folders across orgID",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: orgId2},
			},
			source:          "bravo",
			destination:     "foxtrot",
			want:            nil,
			wantRunError:    true,
		},
		{
			name: "11. Invalid source folder name",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "delta", Paths: "alpha.delta", OrgId: orgId1},
			},
			source: "invalid_folder",
			destination: "delta",
			want: nil,
			wantRunError: true,
		},
		{
			name: "12. Invalid destination folder name",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
			},
			source: "bravo",
			destination: "invalid_folder",
			want: nil,
			wantRunError: true,
		},
		{
			name: "13. Invalid destination AND source folder names",
			initialFolders: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: orgId1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: orgId1},
			},
			source: "invalid1",
			destination: "invalid2",
			want: nil,
			wantRunError: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			driver, err := folder.NewDriver(tt.initialFolders)
			if err != nil {
				t.Fatalf("driver failed to initialise - %v", err)
			}

			got, err := driver.MoveFolder(tt.source, tt.destination)

			if tt.wantRunError && err == nil {
				t.Errorf("%v: MoveFolder() expected an error but got none", tt.name)
				return
			}


			if !tt.wantRunError && err != nil {
				t.Errorf("MoveFolder() received unexpected error: %v", err)
				return
			}

			// shallow folder comparison
			if !compareFolders(got, tt.want) {
				t.Errorf("MoveFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}

