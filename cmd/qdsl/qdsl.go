// Copyright 2023 NJWS Inc.
// Copyright 2022 Listware

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"github.com/manifoldco/promptui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var options = qdsl.NewOptions()

var customFilter []string

var confirm bool

func buildFilter(query string, filters []string) string {
	if len(filters) == 0 {
		return query
	}
	// find *. or <. at the start
	r, err := regexp.Compile(`^[\*\<]\.`)
	if err != nil {
		log.Error("Can't add filter, parse error: ", err)
		return query
	}
	num := r.FindStringIndex(query)
	if num != nil {
		i := num[1] - 1
		query = query[:i] + "[?" + filters[0] + "?]" + query[i:]
		filters = filters[1:]
	}

	if strings.Contains(query, "[?") {
		i := strings.Index(query, "?]")
		if i == -1 {
			log.Error("Error while parsing filter: can't find close filter operator '?]'")
			return query
		}
		newQuery := query[:i]
		for _, filter := range filters {
			newQuery += " && " + filter
		}
		newQuery += query[i:]
		query = newQuery
		// return newQuery
	}
	log.WithFields(logrus.Fields{"cli": "qdsl"}).Debug(query)

	return query
}

func qdslQuery(cmd *cobra.Command, args []string) {
	log.WithFields(logrus.Fields{"cli": "qdsl"}).Debug("QDSL called with argument: ", args[0])
	query := buildFilter(args[0], customFilter)

	if options.Remove && !confirm {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Confirm remove %s", query),
			IsConfirm: true,
		}

		_, err := prompt.Run()
		if err != nil {
			return
		}
	}

	elements, err := qdsl.RawQdsl(context.Background(), query, options)
	if err != nil {
		log.Error(err)
		return
	}

	s, _ := json.Marshal(elements)
	fmt.Println(string(s))
}
