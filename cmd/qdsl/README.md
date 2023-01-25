# Query Database Search Language
* search and remove objects/links by qdsl query

###### Example
```
qdsl *.root
qdsl --remove --id *.root
qdsl --remove --linkid *.root
```
###### Flags
* `id - get vertex '_id', default true` 
* `key - get vertex '_key', default false` 
* `object - get vertex 'object', default false` 
* `link - get edge 'object', default false` 
* `linkId - get edge 'id', default false` 
* `name - get edge '_name', default false` 
* `type - get edge '_type', default false` 
* `remove - remove all results, default false` 
* `confirm - confirm remove, default false` 