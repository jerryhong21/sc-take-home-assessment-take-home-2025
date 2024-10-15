package folder_test

import (
	"sort"
	"github.com/georgechieng-sc/interns-2022/folder"
)

// this is a shallow comparison of folders that ignores parent and children comparison
func compareFolders(got, want []folder.Folder) bool {

    if len(got) != len(want) {
        return false
    }

	// sort them
	sortFolders(got, lessThanComparator)
	sortFolders(want, lessThanComparator)

    for i := range got {
        if got[i].Name != want[i].Name || got[i].Paths != want[i].Paths || got[i].OrgId != want[i].OrgId {
            return false
        }
    }

    return true
}

func sortFolders(folders []folder.Folder, lessThanComparator func(f1, f2 folder.Folder) bool) {
	sort.Slice(folders, func(i, j int) bool {
		return lessThanComparator(folders[i], folders[j])
	})
}

// name and orgid is primary key
func lessThanComparator(f1, f2 folder.Folder) bool {
	if f1.Name != f2.Name {
		return f1.Name < f2.Name
	}
	return f1.OrgId.String() < f2.OrgId.String()
}
