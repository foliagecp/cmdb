// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package qdsl

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"git.fg-tech.ru/listware/proto/sdk/pbcmdb/pbqdsl"
)

var (
	// pathLinks = []string{"links"}
	rootName = "root"
)

type refSpec struct {
	specialLevels, queryLevel int
	hasCatchall               bool
}
type refPath struct {
	v, e string
}

func omitRoot(qdsl Path) Path {
	rootBlock := &Block{Node: &Node{Name: &rootName}}

	max := len(qdsl) - 1
	if max == -1 {
		return append(qdsl, rootBlock)
	}

	last := qdsl[max]
	if last.Node == nil {
		return append(qdsl, rootBlock)
	}

	if last.Node.Name == nil {
		return append(qdsl, rootBlock)
	}

	if *last.Node.Name != rootName {
		return append(qdsl, rootBlock)
	}

	return qdsl
}

func pathToAql(element *Element, options *pbqdsl.Options) {
	path := omitRoot(element.Path)

	max := len(path) - 1

	root := path[max]

	qdsl := path[:max]

	length := len(qdsl)

	hasPath := length > 0

	// <...
	hasCatchallValue := 0

	hasCatchall := hasPath && qdsl[0].Catchall
	if hasCatchall {
		hasCatchallValue = 1
	}

	specialLevels := getSpecialDepth(options)

	queryLevel := int(math.Max(float64(length-hasCatchallValue-specialLevels), 0))
	queryLevelMax := queryLevel

	qdslBase := make(Path, len(qdsl))
	copy(qdslBase, qdsl)

	if hasCatchall {
		queryLevelMax = int(math.Max(float64(queryLevel), 10))

		qdslBase := make(Path, len(qdsl)-1)
		copy(qdslBase, qdsl[1:])
	}

	var aqlFraments []string

	// объект, линк, путь на уровень 1..1
	// начальный объект
	// коллекиця линков и/или граф
	aqlFraments = append(aqlFraments, fmt.Sprintf("for object, link, path in %d..%d outbound 'system/%s' graph system\n", queryLevel, queryLevelMax, *root.Node.Name))

	// непонятная магия, которая вроде бы не влияет на результат, но добавляет задержку
	// aqlFraments = append(aqlFraments, autoRestrict("path", false))

	// protect against missing object
	aqlFraments = append(aqlFraments, "filter object\n")

	reverse(qdslBase)

	for i, level := range qdslBase {
		res := convertLevelToAQL(level, i, refSpec{specialLevels, queryLevel, hasCatchall})
		aqlFraments = append(aqlFraments, res)
	}

	if hasCatchall {
		res := convertLevelToAQL(qdsl[0], -1, refSpec{specialLevels, queryLevel, hasCatchall})
		aqlFraments = append(aqlFraments, res)
	}

	querySearch := assembleFrags(aqlFraments, "\n")

	var returnList []string

	if options.Object {
		returnList = append(returnList, "object: object")
	}

	if options.Id {
		returnList = append(returnList, "id: object._id")
	}

	if options.Key {
		returnList = append(returnList, "key: object._key")
	}

	if options.Link {
		returnList = append(returnList, "link: link")
	}

	if options.LinkId {
		returnList = append(returnList, "link_id: link._id")
	}

	if options.Name {
		returnList = append(returnList, "name: link._name")
	}

	if options.Type {
		returnList = append(returnList, "type: link._type")
	}

	if options.Path {
		returnList = append(returnList, "path: path")
	}

	returnExpr := assembleFrags(returnList, ",\n  ")
	element.Query = fmt.Sprintf(`
%s

return {
  %s
}
    `, querySearch, returnExpr)
}

func getSpecialDepth(options *pbqdsl.Options) int {
	return 0
}

// func autoRestrict(pathVar string, reverse bool) string {
// 	var reverseStr = "iv"
// 	if reverse {
// 		reverseStr = "0, iv"
// 	}
// 	return fmt.Sprintf(
// 		`
// filter (for iv in 0..(length(%s.vertices) - 1)
//           let r = %s.vertices[iv]._meta.restrict
//           return iv < 0 || !r || (for e in slice(%s.edges, %s)
//                                   return parse_identifier(e).collection) all in r
//        ) all == true`, pathVar, pathVar, pathVar, reverseStr)
// }

func assembleFrags(frags []string, separator string) string {
	// FIXME remove '\n' from array
	return strings.Join(frags, separator)
}

func convertLevelToAQL(block *Block, i int, refSpec refSpec) string {
	refs := getRefPath(i, refSpec)

	var nameFilter string

	if block.Node != nil {
		names := enrollNodeName(block.Node)
		var condition string
		if len(names) == 1 {
			condition = fmt.Sprintf("== '%s'", names[0])
		} else {
			condition = fmt.Sprintf("in [%s]", strings.Join(names, ","))
		}
		nameFilter = fmt.Sprintf("filter %s._name %s\n", refs.e, condition)
	}
	var attrFilter []string
	if block.Filter != nil {
		for _, filter := range block.Filter.Filter {
			var items []string
			for _, exp := range filter {
				prefix := exp.Expression.Variable[0]
				propPath := exp.Expression.Variable[1:]

				var vertexPivot bool
				if prefix == "@" || prefix == "object" {
					vertexPivot = true
				}

				var edgePivot bool
				if prefix == "$" || prefix == "link" {
					edgePivot = true
				}

				pivot := prefix
				if vertexPivot {
					pivot = refs.v
				} else if edgePivot {
					pivot = refs.e
				}

				var propAccess string

				for _, value := range propPath {
					propAccess = fmt.Sprintf("%s.%s", propAccess, value)
				}

				var result []string
				result = append(result, fmt.Sprintf("%s%s", pivot, propAccess))

				result = append(result, exp.Expression.Op)
				result = append(result, exp.Expression.Evaluation)
				result = append(result, exp.BoolOp)

				items = append(items, assembleFrags(result, " "))
			}

			attrFilter = append(attrFilter, "filter "+assembleFrags(items, " "))
		}
	}

	return nameFilter + assembleFrags(attrFilter, "\n")
}

func single(pathVar string, i int) refPath {
	j := i + 1
	//if i > 0 {
	//j = i + 1
	//}
	return refPath{
		v: fmt.Sprintf("%s.vertices[%d]", pathVar, j),
		e: fmt.Sprintf("%s.edges[%d]", pathVar, i),
	}
}

func getRefPath(i int, refSpec refSpec) refPath {
	if refSpec.specialLevels == 0 || (i >= 0 && i < refSpec.queryLevel) {
		return single("path", i)
	}

	if !refSpec.hasCatchall && i >= refSpec.queryLevel {
		return single("specialPath", i-refSpec.queryLevel)
	}

	if i < 0 {
		if i < -refSpec.specialLevels {
			return single("path", i+refSpec.specialLevels)
		} else {
			return single("specialPath", i)
		}
	}
	// includeDescendants && specialLevels > 0 && i >= queryLevel
	return refPath{
		v: fmt.Sprintf("(path.vertices[%d] || specialPath.vertices[%d - length(path.vertices)])", i+1, i+1),
		e: fmt.Sprintf("(path.edges[%d] || specialPath.edges[%d - length(path.edges)])", i, i),
	}
}

func enrollNodeName(node *Node) (names []string) {
	if node.Ranges == nil {
		names = append(names, *node.Name)
		return
	}

	for _, ranges := range node.Ranges {
		from := *ranges.From

		if ranges.To == nil {
			if node.Name == nil {
				names = append(names, from)
			} else {
				names = append(names, *node.Name+from)
			}
		} else {
			to := *ranges.To

			var isPadded bool
			var fixedSize int

			if from[0] == '0' {
				isPadded = true
				fixedSize = len(from)
			}

			fromInt, _ := strconv.Atoi(from)
			toInt, _ := strconv.Atoi(to)

			for i := fromInt; i <= toInt; i++ {
				if isPadded {
					names = append(names, fmt.Sprintf("\"%s%s\"", *node.Name, padZeroes(i, fixedSize)))
				} else {
					names = append(names, fmt.Sprintf("\"%s%d\"", *node.Name, i))
				}

			}
		}

	}
	return
}

func padZeroes(i, paddedSize int) (s string) {
	s = fmt.Sprint(i)

	for len(s) < paddedSize {
		s = "0" + s
	}

	return
}

// func getQuerySort(qdsl *Block) (s string) {
// 	if qdsl == nil {
// 		return
// 	}

// 	return
// }

/*
function getQueryLimit(qdsl) {
  const l = (qdsl && qdsl.limits) ? qdsl.limits : null;
  return l
      ? `limit ${l.offset}, ${l.limit}`
      : '';
}

function getQuerySort(qdsl) {
    const sortStr = (!qdsl || !qdsl.sort || qdsl.sort.length === 0)
    ? null
    : qdsl.sort.reduce((str, el) => {
        let currentPrefix;
        const [prefix, ...rest] = el.field;

        switch(prefix) {
        case '@':
        case 'object':
            currentPrefix = 'rObject';
            break;
        case '$':
        case 'link':
            currentPrefix = 'rLink';
            break;
        case 'path':
            currentPrefix = 'rPath';
            break;
        default:
            currentPrefix = prefix;
        }

        const fieldStr = [currentPrefix, ...rest].join('.');

        return str
        ? `${str}, ${fieldStr} ${el.direction}`
        : `sort ${fieldStr} ${el.direction}`;
    }, null);

    return sortStr ? sortStr + '\n' : '';
}*/
