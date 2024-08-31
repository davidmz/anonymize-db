package main

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/davidmz/mustbe"
	lorem "github.com/drhodes/golorem"
)

const EOT = "\\."
const NULL = "\\N"

var headerLineRe = regexp.MustCompile(`^COPY\s+(\S+)\s+\((.+?)\)`)

func processTable(headerLine string, scanner *bufio.Scanner) {
	mustbe.OKVal(fmt.Fprintln(output, headerLine))

	header := headerLineRe.FindStringSubmatch(headerLine)
	mustbe.True(header != nil, fmt.Errorf("invalid COPY line: %s", headerLine))
	//lint:ignore SA5011 we are sure that header is not nil, see prev. line
	table := header[1]
	//lint:ignore SA5011 we are sure that header is not nil here too
	columns := strings.Split(header[2], ", ")
	log.Println("Table found:", table, columns)

	tableCfg := config.Tables[table]
	colRules := make([]string, len(columns))
	for colName, rule := range tableCfg.Columns {
		for i, c := range columns {
			if c == colName {
				colRules[i] = rule
				break
			}
		}
	}

	switch {
	case tableCfg.Keep:
		log.Println("Keep this table")
		for scanner.Scan() {
			line := scanner.Text()
			mustbe.OKVal(fmt.Fprintln(output, line))
			if line == EOT {
				return
			}
		}
		mustbe.OK(scanner.Err())
	case tableCfg.Clean:
		log.Println("Clean this table")
		for scanner.Scan() {
			line := scanner.Text()
			if line == EOT {
				mustbe.OKVal(fmt.Fprintln(output, line))
				return
			}
		}
		mustbe.OK(scanner.Err())
	default:
		log.Println("Anonymize this table")
		for scanner.Scan() {
			line := scanner.Text()
			if line == EOT {
				mustbe.OKVal(fmt.Fprintln(output, line))
				return
			}

			if config.EncryptUUIDs {
				// Anonymize all UUIDs in the line
				line = anonAllUUIDs(line)
			}
			if len(config.Tables[table].Columns) > 0 {
				values := strings.Split(line, "\t")
				for i, rule := range colRules {
					if values[i] != "" && values[i] != NULL {
						switch {
						case rule == "shorttext":
							values[i] = strings.TrimSuffix(lorem.Sentence(1, 3), ".")
						case rule == "text":
							values[i] = lorem.Paragraph(2, 5)
						case rule == "uniqword":
							values[i] = anonWord(values[i])
						case rule == "uniqemail":
							values[i] = anonEmail(values[i])
						case rule == "shortid":
							values[i] = anonShortId(values[i])
						case rule == "uuids" && !config.EncryptUUIDs:
							values[i] = anonAllUUIDs(values[i])
						case strings.HasPrefix(rule, "set:"):
							values[i] = strings.TrimPrefix(rule, "set:")
						}
					}
				}
				line = strings.Join(values, "\t")
			}
			mustbe.OKVal(fmt.Fprintln(output, line))
		}
		mustbe.OK(scanner.Err())
	}
	mustbe.Thrown(fmt.Errorf("premature end of file"))
}
