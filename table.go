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
	table := header[1]
	columns := strings.Split(header[2], ", ")
	log.Println("Table found:", table, columns)

	colRules := make([]string, len(columns))
	for colName, rule := range rules[table].Columns {
		for i, c := range columns {
			if c == colName {
				colRules[i] = rule
				break
			}
		}
	}

	switch rules[table].Action {
	case "KEEP":
		log.Println("Keep this table")
		for scanner.Scan() {
			line := scanner.Text()
			mustbe.OKVal(fmt.Fprintln(output, line))
			if line == EOT {
				return
			}
		}
		mustbe.OK(scanner.Err())
	case "CLEAN":
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

			// Anonymize all UUIDs in the line
			line = anonUUIDs(line)
			if len(rules[table].Columns) > 0 {
				values := strings.Split(line, "\t")
				for i, rule := range colRules {
					if values[i] != "" && values[i] != NULL {
						switch rule {
						case "shorttext":
							values[i] = strings.TrimSuffix(lorem.Sentence(1, 3), ".")
						case "text":
							values[i] = lorem.Paragraph(2, 5)
						case "uniqword":
							values[i] = anonWord(values[i])
						case "uniqemail":
							values[i] = anonEmail(values[i])
						default:
							if strings.HasPrefix(rule, "set:") {
								values[i] = strings.TrimPrefix(rule, "set:")
							}
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
