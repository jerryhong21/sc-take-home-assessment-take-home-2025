## Assumptions in implementation adn testing

1. Every folder will have a primary key (orgId, name) - otherwise getAllChildrenFolder will be ambiguous since it only takes in orgId and folder name
2. There can only be one orgId that contains a specific pair of folders with the names "name1" and "name2". 
    - E.g. if the folders "alpha" and "beta" exist in orgId1, no other orgId can have folders with the same names "alpha" and "beta" as a pair.