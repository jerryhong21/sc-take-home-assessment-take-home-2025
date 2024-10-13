


### Blog and thoughts


#### Initial thoughts and learning GoLang \(13/10/24\) 
- The data structure is a Folder slice, this means O(n) opereation for both move folder and getAllChild folder since a linear traversal is O(n) where `n` is the number of folders.
- Goroutines and channels are currently used in data generation in `static.go`
- Thinking of pivoting to an alternate data structure, with consideration of other functionalities such as:
    - deleting folder
    - creating folder
    - retrieve folder path (pwd)
    - modifying folder (rename)
    - searching folder (by name, by path, by orgId etc)

- All of these functionalities, as well as the natural folder hierarchy inspires a *tree-based structure with parent and child references*.

- This optimise search time in the average case, whilst worst case time complexity remains O(n);
- Additionally, supplementary maps can be used for quick lookups 
    - maps `fullpath` to folder pointers, this is worst case O(1).
    - maps `folderName` to slices of folder pointers - O(1) average
    - maps `orgId` to slices of folder pointers (this is questionable since there contains many duplcaites of the same orgId) - this is O(1) average

``` go
type Folder struct {
    Name     string
    OrgID    uuid.UUID
    Path     string
    Parent   *Folder
    Children []*Folder
}
```

Using this data strucutre,

`GetAllChildFolders`:


Has direct access to children slice - $O(1)$ retrial time. retrieving all descendants is $O(k)$ Where $k$ is the number of descendant folders. This is on aveerage more efficient.

`MoveFolder`:

Updating parent pointers and adjusting path attributes of all original descendents can be managed recursively. This is $O(k)$.


Space complexity is also reduced to $O(k)$ since we just need to temporarily store the descendent paths.

#### Tradeoffs discussion
The tradeoffs of using this data strucutre are apparent - the benefit is efficiency, especially in consideration for other operations whilst the downside is added complexity, and implementation maintenance to ensure data structures are updated appropriately.