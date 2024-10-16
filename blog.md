


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

`GetAllChildFolders`:


Has direct access to children slice - $O(1)$ retrial time. retrieving all descendants is $O(k)$ Where $k$ is the number of descendant folders. This is on aveerage more efficient.

`MoveFolder`:

Updating parent pointers and adjusting path attributes of all original descendents can be managed recursively. This is $O(k)$.


Space complexity is also reduced to $O(k)$ since we just need to temporarily store the descendent paths.

#### Tradeoffs discussion
The tradeoffs of using this data strucutre are apparent - the benefit is efficiency, especially in consideration for other operations whilst the downside is added complexity, and implementation maintenance to ensure data structures are updated appropriately.

Another major benefit is scalability, hierarchical relationships will be able to perform better with deep hierarchies and a larger database as opposed to a flat slice.



#### Test planning for `getAllChildFolders` function  \(15/10/24\) 


##### Possible test cases:

##### Basic functionalities include
1. Single node (no children)
2. Parent with multiple children
3. Root node - no parent exists
4. Nested child folders with depth

###### Edge cases, these most likely deal with erraneous look up operations, duplicate folder names, duplicate orgids or other attributes

- Non-existent parents
- Parent in a different OrgID
- Invalid orgId
- Empty / invalid parent name (this can include non-ascii characters which are important of internationalisation)
- Case sensitivity - difference cases of folder names
- Parent is a leaf node 
- check malformed folder paths
- parent name as a substring of exisitng folder (avoid partial matches)
- Check for cyclic paths






#### Test planning for `MoveFolder` function  \(16/10/24\) 

#### Functionalities
1. moving folder with no children into another folder with no children
2. moving folder with no children into another folder with children
3. moving folder with children into another folder with no children
4. moving folder with children into another folder with children
5. Move a deeply nested folder into another folder with children
6. move folder in-place (move to its current parent)
7. move folder when source is root folder


##### Edge cases and error checking

1. Moving a node to a child of itself
2. Moving a folder to itself
3. Moving folders across orgId
4. Invalid source folder name
5. Invalid dest folder name







