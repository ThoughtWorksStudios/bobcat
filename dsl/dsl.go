package dsl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

func invalid(format string, tokens ...interface{}) error {
	return fmt.Errorf(format, tokens...)
}

// print arbitrary messages to STDERR; useful when making debug statements
// for development
func debug(format string, tokens ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", tokens...)
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Script",
			pos:  position{line: 15, col: 1, offset: 335},
			expr: &actionExpr{
				pos: position{line: 15, col: 10, offset: 344},
				run: (*parser).callonScript1,
				expr: &seqExpr{
					pos: position{line: 15, col: 10, offset: 344},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 15, col: 10, offset: 344},
							label: "prog",
							expr: &zeroOrMoreExpr{
								pos: position{line: 15, col: 15, offset: 349},
								expr: &ruleRefExpr{
									pos:  position{line: 15, col: 15, offset: 349},
									name: "Statement",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 15, col: 26, offset: 360},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Statement",
			pos:  position{line: 19, col: 1, offset: 396},
			expr: &actionExpr{
				pos: position{line: 19, col: 13, offset: 408},
				run: (*parser).callonStatement1,
				expr: &labeledExpr{
					pos:   position{line: 19, col: 13, offset: 408},
					label: "statement",
					expr: &choiceExpr{
						pos: position{line: 19, col: 24, offset: 419},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 19, col: 24, offset: 419},
								name: "EntityExpr",
							},
							&ruleRefExpr{
								pos:  position{line: 19, col: 37, offset: 432},
								name: "GenerateOverrideExpr",
							},
							&ruleRefExpr{
								pos:  position{line: 19, col: 60, offset: 455},
								name: "GenerateExpr",
							},
						},
					},
				},
			},
		},
		{
			name: "GenerateExpr",
			pos:  position{line: 23, col: 1, offset: 498},
			expr: &actionExpr{
				pos: position{line: 23, col: 16, offset: 513},
				run: (*parser).callonGenerateExpr1,
				expr: &seqExpr{
					pos: position{line: 23, col: 16, offset: 513},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 23, col: 16, offset: 513},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 23, col: 18, offset: 515},
							val:        "generate",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 23, col: 29, offset: 526},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 23, col: 31, offset: 528},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 23, col: 36, offset: 533},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 23, col: 47, offset: 544},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 23, col: 49, offset: 546},
							label: "args",
							expr: &ruleRefExpr{
								pos:  position{line: 23, col: 54, offset: 551},
								name: "Arguments",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 23, col: 64, offset: 561},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "GenerateOverrideExpr",
			pos:  position{line: 27, col: 1, offset: 605},
			expr: &actionExpr{
				pos: position{line: 27, col: 24, offset: 628},
				run: (*parser).callonGenerateOverrideExpr1,
				expr: &seqExpr{
					pos: position{line: 27, col: 24, offset: 628},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 27, col: 24, offset: 628},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 27, col: 26, offset: 630},
							val:        "generate",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 27, col: 37, offset: 641},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 27, col: 39, offset: 643},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 27, col: 44, offset: 648},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 27, col: 55, offset: 659},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 27, col: 57, offset: 661},
							label: "args",
							expr: &ruleRefExpr{
								pos:  position{line: 27, col: 62, offset: 666},
								name: "Arguments",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 27, col: 72, offset: 676},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 27, col: 74, offset: 678},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 27, col: 78, offset: 682},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 27, col: 80, offset: 684},
							label: "body",
							expr: &zeroOrOneExpr{
								pos: position{line: 27, col: 85, offset: 689},
								expr: &ruleRefExpr{
									pos:  position{line: 27, col: 85, offset: 689},
									name: "FieldSet",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 27, col: 95, offset: 699},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 27, col: 97, offset: 701},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 27, col: 101, offset: 705},
							name: "_",
						},
					},
				},
			},
		},
		{
			name:        "EntityExpr",
			displayName: "\"entity declaration\"",
			pos:         position{line: 31, col: 1, offset: 750},
			expr: &choiceExpr{
				pos: position{line: 31, col: 35, offset: 784},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 31, col: 35, offset: 784},
						run: (*parser).callonEntityExpr2,
						expr: &seqExpr{
							pos: position{line: 31, col: 35, offset: 784},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 31, col: 35, offset: 784},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 31, col: 37, offset: 786},
									val:        "def",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 31, col: 43, offset: 792},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 31, col: 45, offset: 794},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 31, col: 50, offset: 799},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 31, col: 61, offset: 810},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 31, col: 63, offset: 812},
									val:        "{",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 31, col: 67, offset: 816},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 31, col: 69, offset: 818},
									label: "body",
									expr: &zeroOrOneExpr{
										pos: position{line: 31, col: 74, offset: 823},
										expr: &ruleRefExpr{
											pos:  position{line: 31, col: 74, offset: 823},
											name: "FieldSet",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 31, col: 84, offset: 833},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 31, col: 86, offset: 835},
									val:        "}",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 31, col: 90, offset: 839},
									name: "_",
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 33, col: 5, offset: 882},
						name: "FailOnUnterminatedEntity",
					},
				},
			},
		},
		{
			name:        "FieldSet",
			displayName: "\"entity fields\"",
			pos:         position{line: 35, col: 1, offset: 908},
			expr: &choiceExpr{
				pos: position{line: 35, col: 28, offset: 935},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 35, col: 28, offset: 935},
						name: "FailOnUndelimitedFields",
					},
					&actionExpr{
						pos: position{line: 35, col: 54, offset: 961},
						run: (*parser).callonFieldSet3,
						expr: &seqExpr{
							pos: position{line: 35, col: 54, offset: 961},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 35, col: 54, offset: 961},
									label: "first",
									expr: &ruleRefExpr{
										pos:  position{line: 35, col: 60, offset: 967},
										name: "FieldDecl",
									},
								},
								&labeledExpr{
									pos:   position{line: 35, col: 70, offset: 977},
									label: "rest",
									expr: &zeroOrMoreExpr{
										pos: position{line: 35, col: 75, offset: 982},
										expr: &seqExpr{
											pos: position{line: 35, col: 76, offset: 983},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 35, col: 76, offset: 983},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 35, col: 78, offset: 985},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 35, col: 82, offset: 989},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 35, col: 84, offset: 991},
													name: "FieldDecl",
												},
											},
										},
									},
								},
								&zeroOrOneExpr{
									pos: position{line: 35, col: 96, offset: 1003},
									expr: &seqExpr{
										pos: position{line: 35, col: 97, offset: 1004},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 35, col: 97, offset: 1004},
												name: "_",
											},
											&litMatcher{
												pos:        position{line: 35, col: 99, offset: 1006},
												val:        ",",
												ignoreCase: false,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "FieldDecl",
			pos:  position{line: 39, col: 1, offset: 1063},
			expr: &choiceExpr{
				pos: position{line: 39, col: 13, offset: 1075},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 39, col: 13, offset: 1075},
						name: "StaticDecl",
					},
					&ruleRefExpr{
						pos:  position{line: 39, col: 26, offset: 1088},
						name: "SymbolicDecl",
					},
				},
			},
		},
		{
			name:        "StaticDecl",
			displayName: "\"field declaration\"",
			pos:         position{line: 41, col: 1, offset: 1102},
			expr: &actionExpr{
				pos: position{line: 41, col: 34, offset: 1135},
				run: (*parser).callonStaticDecl1,
				expr: &seqExpr{
					pos: position{line: 41, col: 34, offset: 1135},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 41, col: 34, offset: 1135},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 41, col: 39, offset: 1140},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 41, col: 50, offset: 1151},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 41, col: 52, offset: 1153},
							label: "fieldValue",
							expr: &ruleRefExpr{
								pos:  position{line: 41, col: 63, offset: 1164},
								name: "Literal",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 41, col: 71, offset: 1172},
							name: "_",
						},
					},
				},
			},
		},
		{
			name:        "SymbolicDecl",
			displayName: "\"field declaration\"",
			pos:         position{line: 45, col: 1, offset: 1225},
			expr: &actionExpr{
				pos: position{line: 45, col: 36, offset: 1260},
				run: (*parser).callonSymbolicDecl1,
				expr: &seqExpr{
					pos: position{line: 45, col: 36, offset: 1260},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 45, col: 36, offset: 1260},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 45, col: 41, offset: 1265},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 45, col: 52, offset: 1276},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 45, col: 54, offset: 1278},
							label: "fieldType",
							expr: &ruleRefExpr{
								pos:  position{line: 45, col: 64, offset: 1288},
								name: "Builtin",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 45, col: 72, offset: 1296},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 45, col: 74, offset: 1298},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 45, col: 79, offset: 1303},
								expr: &ruleRefExpr{
									pos:  position{line: 45, col: 79, offset: 1303},
									name: "Arguments",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 45, col: 90, offset: 1314},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "Arguments",
			pos:  position{line: 49, col: 1, offset: 1373},
			expr: &choiceExpr{
				pos: position{line: 49, col: 13, offset: 1385},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 49, col: 13, offset: 1385},
						run: (*parser).callonArguments2,
						expr: &seqExpr{
							pos: position{line: 49, col: 13, offset: 1385},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 49, col: 13, offset: 1385},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 49, col: 17, offset: 1389},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 49, col: 19, offset: 1391},
									label: "body",
									expr: &zeroOrOneExpr{
										pos: position{line: 49, col: 24, offset: 1396},
										expr: &ruleRefExpr{
											pos:  position{line: 49, col: 24, offset: 1396},
											name: "ArgumentsBody",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 49, col: 39, offset: 1411},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 49, col: 41, offset: 1413},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 51, col: 5, offset: 1463},
						name: "FailOnUnterminatedArguments",
					},
				},
			},
		},
		{
			name:        "ArgumentsBody",
			displayName: "\"arguments body\"",
			pos:         position{line: 53, col: 1, offset: 1492},
			expr: &choiceExpr{
				pos: position{line: 53, col: 34, offset: 1525},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 53, col: 34, offset: 1525},
						name: "FailOnUndelimitedArgs",
					},
					&actionExpr{
						pos: position{line: 53, col: 58, offset: 1549},
						run: (*parser).callonArgumentsBody3,
						expr: &seqExpr{
							pos: position{line: 53, col: 58, offset: 1549},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 53, col: 58, offset: 1549},
									label: "first",
									expr: &ruleRefExpr{
										pos:  position{line: 53, col: 64, offset: 1555},
										name: "SingleArgument",
									},
								},
								&labeledExpr{
									pos:   position{line: 53, col: 79, offset: 1570},
									label: "rest",
									expr: &zeroOrMoreExpr{
										pos: position{line: 53, col: 84, offset: 1575},
										expr: &seqExpr{
											pos: position{line: 53, col: 85, offset: 1576},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 53, col: 85, offset: 1576},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 53, col: 87, offset: 1578},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 53, col: 91, offset: 1582},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 53, col: 93, offset: 1584},
													name: "SingleArgument",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Literal",
			pos:  position{line: 57, col: 1, offset: 1652},
			expr: &choiceExpr{
				pos: position{line: 57, col: 11, offset: 1662},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 57, col: 11, offset: 1662},
						name: "DateTimeLiteral",
					},
					&ruleRefExpr{
						pos:  position{line: 57, col: 29, offset: 1680},
						name: "NumberLiteral",
					},
					&ruleRefExpr{
						pos:  position{line: 57, col: 45, offset: 1696},
						name: "BoolLiteral",
					},
					&ruleRefExpr{
						pos:  position{line: 57, col: 59, offset: 1710},
						name: "StringLiteral",
					},
					&ruleRefExpr{
						pos:  position{line: 57, col: 75, offset: 1726},
						name: "NullLiteral",
					},
				},
			},
		},
		{
			name: "SingleArgument",
			pos:  position{line: 59, col: 1, offset: 1739},
			expr: &choiceExpr{
				pos: position{line: 59, col: 18, offset: 1756},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 59, col: 18, offset: 1756},
						name: "Literal",
					},
					&ruleRefExpr{
						pos:  position{line: 59, col: 28, offset: 1766},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 61, col: 1, offset: 1778},
			expr: &choiceExpr{
				pos: position{line: 61, col: 14, offset: 1791},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 61, col: 14, offset: 1791},
						run: (*parser).callonIdentifier2,
						expr: &seqExpr{
							pos: position{line: 61, col: 14, offset: 1791},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 61, col: 14, offset: 1791},
									expr: &ruleRefExpr{
										pos:  position{line: 61, col: 15, offset: 1792},
										name: "ReservedWord",
									},
								},
								&charClassMatcher{
									pos:        position{line: 61, col: 28, offset: 1805},
									val:        "[a-z_]i",
									chars:      []rune{'_'},
									ranges:     []rune{'a', 'z'},
									ignoreCase: true,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 61, col: 35, offset: 1812},
									expr: &charClassMatcher{
										pos:        position{line: 61, col: 35, offset: 1812},
										val:        "[a-z0-9_]i",
										chars:      []rune{'_'},
										ranges:     []rune{'a', 'z', '0', '9'},
										ignoreCase: true,
										inverted:   false,
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 63, col: 5, offset: 1849},
						name: "FailOnIllegalIdentifier",
					},
				},
			},
		},
		{
			name:        "Builtin",
			displayName: "\"built-in types\"",
			pos:         position{line: 65, col: 1, offset: 1874},
			expr: &actionExpr{
				pos: position{line: 65, col: 28, offset: 1901},
				run: (*parser).callonBuiltin1,
				expr: &ruleRefExpr{
					pos:  position{line: 65, col: 28, offset: 1901},
					name: "FieldTypes",
				},
			},
		},
		{
			name: "DateTimeLiteral",
			pos:  position{line: 69, col: 1, offset: 1941},
			expr: &choiceExpr{
				pos: position{line: 69, col: 19, offset: 1959},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 69, col: 19, offset: 1959},
						run: (*parser).callonDateTimeLiteral2,
						expr: &seqExpr{
							pos: position{line: 69, col: 19, offset: 1959},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 69, col: 19, offset: 1959},
									label: "date",
									expr: &ruleRefExpr{
										pos:  position{line: 69, col: 24, offset: 1964},
										name: "IsoDate",
									},
								},
								&labeledExpr{
									pos:   position{line: 69, col: 32, offset: 1972},
									label: "localTime",
									expr: &zeroOrOneExpr{
										pos: position{line: 69, col: 42, offset: 1982},
										expr: &ruleRefExpr{
											pos:  position{line: 69, col: 42, offset: 1982},
											name: "LocalTimePart",
										},
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 71, col: 5, offset: 2048},
						name: "FailOnMissingDate",
					},
				},
			},
		},
		{
			name: "LocalTimePart",
			pos:  position{line: 73, col: 1, offset: 2067},
			expr: &actionExpr{
				pos: position{line: 73, col: 17, offset: 2083},
				run: (*parser).callonLocalTimePart1,
				expr: &seqExpr{
					pos: position{line: 73, col: 17, offset: 2083},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 73, col: 17, offset: 2083},
							label: "ts",
							expr: &ruleRefExpr{
								pos:  position{line: 73, col: 20, offset: 2086},
								name: "TimePart",
							},
						},
						&labeledExpr{
							pos:   position{line: 73, col: 29, offset: 2095},
							label: "zone",
							expr: &zeroOrOneExpr{
								pos: position{line: 73, col: 34, offset: 2100},
								expr: &ruleRefExpr{
									pos:  position{line: 73, col: 34, offset: 2100},
									name: "ZonePart",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "IsoDate",
			pos:  position{line: 81, col: 1, offset: 2240},
			expr: &actionExpr{
				pos: position{line: 81, col: 11, offset: 2250},
				run: (*parser).callonIsoDate1,
				expr: &seqExpr{
					pos: position{line: 81, col: 11, offset: 2250},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 81, col: 11, offset: 2250},
							name: "DIGIT",
						},
						&ruleRefExpr{
							pos:  position{line: 81, col: 17, offset: 2256},
							name: "DIGIT",
						},
						&ruleRefExpr{
							pos:  position{line: 81, col: 23, offset: 2262},
							name: "DIGIT",
						},
						&ruleRefExpr{
							pos:  position{line: 81, col: 29, offset: 2268},
							name: "DIGIT",
						},
						&litMatcher{
							pos:        position{line: 81, col: 35, offset: 2274},
							val:        "-",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 81, col: 39, offset: 2278},
							name: "DIGIT",
						},
						&ruleRefExpr{
							pos:  position{line: 81, col: 45, offset: 2284},
							name: "DIGIT",
						},
						&litMatcher{
							pos:        position{line: 81, col: 51, offset: 2290},
							val:        "-",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 81, col: 55, offset: 2294},
							name: "DIGIT",
						},
						&ruleRefExpr{
							pos:  position{line: 81, col: 61, offset: 2300},
							name: "DIGIT",
						},
					},
				},
			},
		},
		{
			name: "TimePart",
			pos:  position{line: 82, col: 1, offset: 2336},
			expr: &actionExpr{
				pos: position{line: 82, col: 12, offset: 2347},
				run: (*parser).callonTimePart1,
				expr: &seqExpr{
					pos: position{line: 82, col: 12, offset: 2347},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 82, col: 12, offset: 2347},
							val:        "t",
							ignoreCase: true,
						},
						&ruleRefExpr{
							pos:  position{line: 82, col: 17, offset: 2352},
							name: "DIGIT",
						},
						&ruleRefExpr{
							pos:  position{line: 82, col: 23, offset: 2358},
							name: "DIGIT",
						},
						&litMatcher{
							pos:        position{line: 82, col: 29, offset: 2364},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 82, col: 33, offset: 2368},
							name: "DIGIT",
						},
						&ruleRefExpr{
							pos:  position{line: 82, col: 39, offset: 2374},
							name: "DIGIT",
						},
						&litMatcher{
							pos:        position{line: 82, col: 45, offset: 2380},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 82, col: 49, offset: 2384},
							name: "DIGIT",
						},
						&ruleRefExpr{
							pos:  position{line: 82, col: 55, offset: 2390},
							name: "DIGIT",
						},
					},
				},
			},
		},
		{
			name: "ZonePart",
			pos:  position{line: 83, col: 1, offset: 2444},
			expr: &choiceExpr{
				pos: position{line: 83, col: 12, offset: 2455},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 83, col: 12, offset: 2455},
						run: (*parser).callonZonePart2,
						expr: &litMatcher{
							pos:        position{line: 83, col: 12, offset: 2455},
							val:        "z",
							ignoreCase: true,
						},
					},
					&actionExpr{
						pos: position{line: 83, col: 39, offset: 2482},
						run: (*parser).callonZonePart4,
						expr: &seqExpr{
							pos: position{line: 83, col: 39, offset: 2482},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 83, col: 39, offset: 2482},
									val:        "[+-]",
									chars:      []rune{'+', '-'},
									ignoreCase: false,
									inverted:   false,
								},
								&ruleRefExpr{
									pos:  position{line: 83, col: 44, offset: 2487},
									name: "DIGIT",
								},
								&ruleRefExpr{
									pos:  position{line: 83, col: 50, offset: 2493},
									name: "DIGIT",
								},
								&zeroOrOneExpr{
									pos: position{line: 83, col: 56, offset: 2499},
									expr: &litMatcher{
										pos:        position{line: 83, col: 56, offset: 2499},
										val:        ":",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 83, col: 61, offset: 2504},
									name: "DIGIT",
								},
								&ruleRefExpr{
									pos:  position{line: 83, col: 67, offset: 2510},
									name: "DIGIT",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NumberLiteral",
			pos:  position{line: 85, col: 1, offset: 2578},
			expr: &choiceExpr{
				pos: position{line: 85, col: 17, offset: 2594},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 85, col: 17, offset: 2594},
						run: (*parser).callonNumberLiteral2,
						expr: &seqExpr{
							pos: position{line: 85, col: 17, offset: 2594},
							exprs: []interface{}{
								&zeroOrOneExpr{
									pos: position{line: 85, col: 17, offset: 2594},
									expr: &litMatcher{
										pos:        position{line: 85, col: 17, offset: 2594},
										val:        "-",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 85, col: 22, offset: 2599},
									name: "INT",
								},
								&zeroOrOneExpr{
									pos: position{line: 85, col: 26, offset: 2603},
									expr: &seqExpr{
										pos: position{line: 85, col: 27, offset: 2604},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 85, col: 27, offset: 2604},
												val:        ".",
												ignoreCase: false,
											},
											&oneOrMoreExpr{
												pos: position{line: 85, col: 31, offset: 2608},
												expr: &ruleRefExpr{
													pos:  position{line: 85, col: 31, offset: 2608},
													name: "DIGIT",
												},
											},
										},
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 91, col: 5, offset: 2760},
						name: "FailOnOctal",
					},
				},
			},
		},
		{
			name: "BoolLiteral",
			pos:  position{line: 93, col: 1, offset: 2773},
			expr: &actionExpr{
				pos: position{line: 93, col: 15, offset: 2787},
				run: (*parser).callonBoolLiteral1,
				expr: &ruleRefExpr{
					pos:  position{line: 93, col: 15, offset: 2787},
					name: "BoolToken",
				},
			},
		},
		{
			name: "NullLiteral",
			pos:  position{line: 97, col: 1, offset: 2830},
			expr: &actionExpr{
				pos: position{line: 97, col: 15, offset: 2844},
				run: (*parser).callonNullLiteral1,
				expr: &ruleRefExpr{
					pos:  position{line: 97, col: 15, offset: 2844},
					name: "NullToken",
				},
			},
		},
		{
			name: "StringLiteral",
			pos:  position{line: 101, col: 1, offset: 2887},
			expr: &actionExpr{
				pos: position{line: 101, col: 17, offset: 2903},
				run: (*parser).callonStringLiteral1,
				expr: &seqExpr{
					pos: position{line: 101, col: 17, offset: 2903},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 101, col: 17, offset: 2903},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 101, col: 21, offset: 2907},
							expr: &choiceExpr{
								pos: position{line: 101, col: 23, offset: 2909},
								alternatives: []interface{}{
									&seqExpr{
										pos: position{line: 101, col: 23, offset: 2909},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 101, col: 23, offset: 2909},
												expr: &ruleRefExpr{
													pos:  position{line: 101, col: 24, offset: 2910},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 101, col: 36, offset: 2922,
											},
										},
									},
									&seqExpr{
										pos: position{line: 101, col: 40, offset: 2926},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 101, col: 40, offset: 2926},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 101, col: 45, offset: 2931},
												name: "EscapeSequence",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 101, col: 63, offset: 2949},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 105, col: 1, offset: 2985},
			expr: &charClassMatcher{
				pos:        position{line: 105, col: 15, offset: 2999},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 107, col: 1, offset: 3015},
			expr: &choiceExpr{
				pos: position{line: 107, col: 18, offset: 3032},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 107, col: 18, offset: 3032},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 107, col: 37, offset: 3051},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 109, col: 1, offset: 3066},
			expr: &charClassMatcher{
				pos:        position{line: 109, col: 20, offset: 3085},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 111, col: 1, offset: 3098},
			expr: &seqExpr{
				pos: position{line: 111, col: 17, offset: 3114},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 111, col: 17, offset: 3114},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 111, col: 21, offset: 3118},
						name: "HEXDIG",
					},
					&ruleRefExpr{
						pos:  position{line: 111, col: 28, offset: 3125},
						name: "HEXDIG",
					},
					&ruleRefExpr{
						pos:  position{line: 111, col: 35, offset: 3132},
						name: "HEXDIG",
					},
					&ruleRefExpr{
						pos:  position{line: 111, col: 42, offset: 3139},
						name: "HEXDIG",
					},
				},
			},
		},
		{
			name: "INT",
			pos:  position{line: 113, col: 1, offset: 3147},
			expr: &choiceExpr{
				pos: position{line: 113, col: 7, offset: 3153},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 113, col: 7, offset: 3153},
						val:        "0",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 113, col: 13, offset: 3159},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 113, col: 13, offset: 3159},
								name: "NON_ZERO",
							},
							&zeroOrMoreExpr{
								pos: position{line: 113, col: 22, offset: 3168},
								expr: &ruleRefExpr{
									pos:  position{line: 113, col: 22, offset: 3168},
									name: "DIGIT",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NON_ZERO",
			pos:  position{line: 115, col: 1, offset: 3176},
			expr: &charClassMatcher{
				pos:        position{line: 115, col: 12, offset: 3187},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "DIGIT",
			pos:  position{line: 117, col: 1, offset: 3194},
			expr: &charClassMatcher{
				pos:        position{line: 117, col: 9, offset: 3202},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HEXDIG",
			pos:  position{line: 119, col: 1, offset: 3209},
			expr: &charClassMatcher{
				pos:        position{line: 119, col: 10, offset: 3218},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "ReservedWord",
			pos:  position{line: 121, col: 1, offset: 3229},
			expr: &choiceExpr{
				pos: position{line: 121, col: 16, offset: 3244},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 121, col: 16, offset: 3244},
						name: "Keyword",
					},
					&ruleRefExpr{
						pos:  position{line: 121, col: 26, offset: 3254},
						name: "FieldTypes",
					},
					&ruleRefExpr{
						pos:  position{line: 121, col: 39, offset: 3267},
						name: "NullToken",
					},
					&ruleRefExpr{
						pos:  position{line: 121, col: 51, offset: 3279},
						name: "BoolToken",
					},
				},
			},
		},
		{
			name: "Keyword",
			pos:  position{line: 123, col: 1, offset: 3290},
			expr: &choiceExpr{
				pos: position{line: 123, col: 11, offset: 3300},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 123, col: 11, offset: 3300},
						val:        "def",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 123, col: 19, offset: 3308},
						val:        "generate",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "FieldTypes",
			pos:  position{line: 125, col: 1, offset: 3320},
			expr: &choiceExpr{
				pos: position{line: 125, col: 14, offset: 3333},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 125, col: 14, offset: 3333},
						val:        "integer",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 125, col: 26, offset: 3345},
						val:        "decimal",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 125, col: 38, offset: 3357},
						val:        "string",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 125, col: 49, offset: 3368},
						val:        "date",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 125, col: 58, offset: 3377},
						val:        "dict",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "NullToken",
			pos:  position{line: 127, col: 1, offset: 3385},
			expr: &litMatcher{
				pos:        position{line: 127, col: 13, offset: 3397},
				val:        "null",
				ignoreCase: false,
			},
		},
		{
			name: "BoolToken",
			pos:  position{line: 129, col: 1, offset: 3405},
			expr: &choiceExpr{
				pos: position{line: 129, col: 13, offset: 3417},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 129, col: 13, offset: 3417},
						val:        "true",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 129, col: 22, offset: 3426},
						val:        "false",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name:        "FailOnOctal",
			displayName: "\"octal numbers not supported\"",
			pos:         position{line: 138, col: 1, offset: 3804},
			expr: &actionExpr{
				pos: position{line: 138, col: 45, offset: 3848},
				run: (*parser).callonFailOnOctal1,
				expr: &seqExpr{
					pos: position{line: 138, col: 45, offset: 3848},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 138, col: 45, offset: 3848},
							val:        "\\0",
							ignoreCase: false,
						},
						&oneOrMoreExpr{
							pos: position{line: 138, col: 51, offset: 3854},
							expr: &ruleRefExpr{
								pos:  position{line: 138, col: 51, offset: 3854},
								name: "DIGIT",
							},
						},
					},
				},
			},
		},
		{
			name:        "FailOnUnterminatedEntity",
			displayName: "\"unterminated entity\"",
			pos:         position{line: 139, col: 1, offset: 3925},
			expr: &actionExpr{
				pos: position{line: 139, col: 50, offset: 3974},
				run: (*parser).callonFailOnUnterminatedEntity1,
				expr: &seqExpr{
					pos: position{line: 139, col: 50, offset: 3974},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 139, col: 50, offset: 3974},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 139, col: 52, offset: 3976},
							val:        "def",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 139, col: 58, offset: 3982},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 139, col: 60, offset: 3984},
							name: "Identifier",
						},
						&ruleRefExpr{
							pos:  position{line: 139, col: 71, offset: 3995},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 139, col: 73, offset: 3997},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 139, col: 77, offset: 4001},
							name: "_",
						},
						&zeroOrOneExpr{
							pos: position{line: 139, col: 79, offset: 4003},
							expr: &ruleRefExpr{
								pos:  position{line: 139, col: 79, offset: 4003},
								name: "FieldSet",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 139, col: 89, offset: 4013},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 139, col: 91, offset: 4015},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name:        "FailOnUndelimitedFields",
			displayName: "\"missing field delimiter\"",
			pos:         position{line: 140, col: 1, offset: 4107},
			expr: &choiceExpr{
				pos: position{line: 140, col: 53, offset: 4159},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 140, col: 53, offset: 4159},
						run: (*parser).callonFailOnUndelimitedFields2,
						expr: &seqExpr{
							pos: position{line: 140, col: 53, offset: 4159},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 140, col: 53, offset: 4159},
									name: "FieldDecl",
								},
								&seqExpr{
									pos: position{line: 140, col: 64, offset: 4170},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 140, col: 64, offset: 4170},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 140, col: 66, offset: 4172},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 140, col: 70, offset: 4176},
											name: "_",
										},
									},
								},
								&oneOrMoreExpr{
									pos: position{line: 140, col: 73, offset: 4179},
									expr: &seqExpr{
										pos: position{line: 140, col: 74, offset: 4180},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 140, col: 74, offset: 4180},
												name: "_",
											},
											&litMatcher{
												pos:        position{line: 140, col: 76, offset: 4182},
												val:        ",",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 140, col: 80, offset: 4186},
												name: "_",
											},
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 140, col: 146, offset: 4252},
						run: (*parser).callonFailOnUndelimitedFields14,
						expr: &seqExpr{
							pos: position{line: 140, col: 146, offset: 4252},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 140, col: 146, offset: 4252},
									name: "FieldDecl",
								},
								&oneOrMoreExpr{
									pos: position{line: 140, col: 156, offset: 4262},
									expr: &seqExpr{
										pos: position{line: 140, col: 157, offset: 4263},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 140, col: 157, offset: 4263},
												name: "_",
											},
											&ruleRefExpr{
												pos:  position{line: 140, col: 159, offset: 4265},
												name: "FieldDecl",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "FailOnUnterminatedArguments",
			displayName: "\"unterminated arguments\"",
			pos:         position{line: 141, col: 1, offset: 4363},
			expr: &actionExpr{
				pos: position{line: 141, col: 56, offset: 4418},
				run: (*parser).callonFailOnUnterminatedArguments1,
				expr: &seqExpr{
					pos: position{line: 141, col: 56, offset: 4418},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 141, col: 56, offset: 4418},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 141, col: 60, offset: 4422},
							name: "_",
						},
						&zeroOrOneExpr{
							pos: position{line: 141, col: 62, offset: 4424},
							expr: &ruleRefExpr{
								pos:  position{line: 141, col: 62, offset: 4424},
								name: "ArgumentsBody",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 141, col: 77, offset: 4439},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 141, col: 80, offset: 4442},
							alternatives: []interface{}{
								&seqExpr{
									pos: position{line: 141, col: 80, offset: 4442},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 141, col: 80, offset: 4442},
											expr: &ruleRefExpr{
												pos:  position{line: 141, col: 81, offset: 4443},
												name: "SingleArgument",
											},
										},
										&charClassMatcher{
											pos:        position{line: 141, col: 96, offset: 4458},
											val:        "[^)]",
											chars:      []rune{')'},
											ignoreCase: false,
											inverted:   true,
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 141, col: 103, offset: 4465},
									name: "EOF",
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "FailOnUndelimitedArgs",
			displayName: "\"missing argument delimiter\"",
			pos:         position{line: 142, col: 1, offset: 4554},
			expr: &actionExpr{
				pos: position{line: 142, col: 54, offset: 4607},
				run: (*parser).callonFailOnUndelimitedArgs1,
				expr: &seqExpr{
					pos: position{line: 142, col: 54, offset: 4607},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 142, col: 54, offset: 4607},
							name: "SingleArgument",
						},
						&oneOrMoreExpr{
							pos: position{line: 142, col: 69, offset: 4622},
							expr: &seqExpr{
								pos: position{line: 142, col: 70, offset: 4623},
								exprs: []interface{}{
									&choiceExpr{
										pos: position{line: 142, col: 71, offset: 4624},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 142, col: 71, offset: 4624},
												name: "_",
											},
											&seqExpr{
												pos: position{line: 142, col: 75, offset: 4628},
												exprs: []interface{}{
													&ruleRefExpr{
														pos:  position{line: 142, col: 75, offset: 4628},
														name: "_",
													},
													&charClassMatcher{
														pos:        position{line: 142, col: 77, offset: 4630},
														val:        "[^,})]",
														chars:      []rune{',', '}', ')'},
														ignoreCase: false,
														inverted:   true,
													},
													&ruleRefExpr{
														pos:  position{line: 142, col: 84, offset: 4637},
														name: "_",
													},
												},
											},
										},
									},
									&ruleRefExpr{
										pos:  position{line: 142, col: 87, offset: 4640},
										name: "SingleArgument",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "FailOnIllegalIdentifier",
			displayName: "\"illegal identifier\"",
			pos:         position{line: 143, col: 1, offset: 4734},
			expr: &actionExpr{
				pos: position{line: 143, col: 48, offset: 4781},
				run: (*parser).callonFailOnIllegalIdentifier1,
				expr: &ruleRefExpr{
					pos:  position{line: 143, col: 48, offset: 4781},
					name: "ReservedWord",
				},
			},
		},
		{
			name:        "FailOnMissingDate",
			displayName: "\"timestamps must have date\"",
			pos:         position{line: 144, col: 1, offset: 4903},
			expr: &actionExpr{
				pos: position{line: 144, col: 49, offset: 4951},
				run: (*parser).callonFailOnMissingDate1,
				expr: &ruleRefExpr{
					pos:  position{line: 144, col: 49, offset: 4951},
					name: "LocalTimePart",
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 153, col: 1, offset: 5168},
			expr: &zeroOrMoreExpr{
				pos: position{line: 153, col: 18, offset: 5185},
				expr: &charClassMatcher{
					pos:        position{line: 153, col: 18, offset: 5185},
					val:        "[ \\t\\r\\n]",
					chars:      []rune{' ', '\t', '\r', '\n'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 155, col: 1, offset: 5197},
			expr: &notExpr{
				pos: position{line: 155, col: 7, offset: 5203},
				expr: &anyMatcher{
					line: 155, col: 8, offset: 5204,
				},
			},
		},
	},
}

func (c *current) onScript1(prog interface{}) (interface{}, error) {
	return rootNode(c, prog)
}

func (p *parser) callonScript1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onScript1(stack["prog"])
}

func (c *current) onStatement1(statement interface{}) (interface{}, error) {
	return statement, nil
}

func (p *parser) callonStatement1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStatement1(stack["statement"])
}

func (c *current) onGenerateExpr1(name, args interface{}) (interface{}, error) {
	return genNode(c, name, nil, args)
}

func (p *parser) callonGenerateExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGenerateExpr1(stack["name"], stack["args"])
}

func (c *current) onGenerateOverrideExpr1(name, args, body interface{}) (interface{}, error) {
	return genNode(c, name, body, args)
}

func (p *parser) callonGenerateOverrideExpr1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGenerateOverrideExpr1(stack["name"], stack["args"], stack["body"])
}

func (c *current) onEntityExpr2(name, body interface{}) (interface{}, error) {
	return entityNode(c, name, body)
}

func (p *parser) callonEntityExpr2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEntityExpr2(stack["name"], stack["body"])
}

func (c *current) onFieldSet3(first, rest interface{}) (interface{}, error) {
	return delimitedNodeSlice(first, rest), nil
}

func (p *parser) callonFieldSet3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldSet3(stack["first"], stack["rest"])
}

func (c *current) onStaticDecl1(name, fieldValue interface{}) (interface{}, error) {
	return staticFieldNode(c, name, fieldValue)
}

func (p *parser) callonStaticDecl1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStaticDecl1(stack["name"], stack["fieldValue"])
}

func (c *current) onSymbolicDecl1(name, fieldType, args interface{}) (interface{}, error) {
	return dynamicFieldNode(c, name, fieldType, args)
}

func (p *parser) callonSymbolicDecl1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSymbolicDecl1(stack["name"], stack["fieldType"], stack["args"])
}

func (c *current) onArguments2(body interface{}) (interface{}, error) {
	return defaultToEmptySlice(body), nil
}

func (p *parser) callonArguments2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArguments2(stack["body"])
}

func (c *current) onArgumentsBody3(first, rest interface{}) (interface{}, error) {
	return delimitedNodeSlice(first, rest), nil
}

func (p *parser) callonArgumentsBody3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArgumentsBody3(stack["first"], stack["rest"])
}

func (c *current) onIdentifier2() (interface{}, error) {
	return idNode(c)
}

func (p *parser) callonIdentifier2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier2()
}

func (c *current) onBuiltin1() (interface{}, error) {
	return builtinNode(c)
}

func (p *parser) callonBuiltin1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBuiltin1()
}

func (c *current) onDateTimeLiteral2(date, localTime interface{}) (interface{}, error) {
	return dateLiteralNode(c, date, localTime)
}

func (p *parser) callonDateTimeLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDateTimeLiteral2(stack["date"], stack["localTime"])
}

func (c *current) onLocalTimePart1(ts, zone interface{}) (interface{}, error) {
	if zone == nil {
		return []string{ts.(string)}, nil
	} else {
		return []string{ts.(string), zone.(string)}, nil
	}
}

func (p *parser) callonLocalTimePart1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLocalTimePart1(stack["ts"], stack["zone"])
}

func (c *current) onIsoDate1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIsoDate1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIsoDate1()
}

func (c *current) onTimePart1() (interface{}, error) {
	return strings.ToUpper(string(c.text)), nil
}

func (p *parser) callonTimePart1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimePart1()
}

func (c *current) onZonePart2() (interface{}, error) {
	return "Z", nil
}

func (p *parser) callonZonePart2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onZonePart2()
}

func (c *current) onZonePart4() (interface{}, error) {
	return strings.Replace(string(c.text), ":", "", -1), nil
}

func (p *parser) callonZonePart4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onZonePart4()
}

func (c *current) onNumberLiteral2() (interface{}, error) {
	if s := string(c.text); strings.ContainsAny(s, ".") {
		return floatLiteralNode(c, s)
	} else {
		return intLiteralNode(c, s)
	}
}

func (p *parser) callonNumberLiteral2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumberLiteral2()
}

func (c *current) onBoolLiteral1() (interface{}, error) {
	return boolLiteralNode(c)
}

func (p *parser) callonBoolLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBoolLiteral1()
}

func (c *current) onNullLiteral1() (interface{}, error) {
	return nullLiteralNode(c)
}

func (p *parser) callonNullLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNullLiteral1()
}

func (c *current) onStringLiteral1() (interface{}, error) {
	return strLiteralNode(c)
}

func (p *parser) callonStringLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringLiteral1()
}

func (c *current) onFailOnOctal1() (interface{}, error) {
	return Node{}, invalid("Octal sequences are not supported")
}

func (p *parser) callonFailOnOctal1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFailOnOctal1()
}

func (c *current) onFailOnUnterminatedEntity1() (interface{}, error) {
	return nil, invalid("Unterminated entity declaration (missing closing curly brace")
}

func (p *parser) callonFailOnUnterminatedEntity1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFailOnUnterminatedEntity1()
}

func (c *current) onFailOnUndelimitedFields2() (interface{}, error) {
	return nil, invalid("Expected another field declaration")
}

func (p *parser) callonFailOnUndelimitedFields2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFailOnUndelimitedFields2()
}

func (c *current) onFailOnUndelimitedFields14() (interface{}, error) {
	return nil, invalid("Multiple field declarations must be delimited with a comma")
}

func (p *parser) callonFailOnUndelimitedFields14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFailOnUndelimitedFields14()
}

func (c *current) onFailOnUnterminatedArguments1() (interface{}, error) {
	return nil, invalid("Unterminated argument list (missing closing parenthesis)")
}

func (p *parser) callonFailOnUnterminatedArguments1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFailOnUnterminatedArguments1()
}

func (c *current) onFailOnUndelimitedArgs1() (interface{}, error) {
	return nil, invalid("Multiple arguments must be delimited with a comma")
}

func (p *parser) callonFailOnUndelimitedArgs1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFailOnUndelimitedArgs1()
}

func (c *current) onFailOnIllegalIdentifier1() (interface{}, error) {
	return Node{Value: string(c.text)}, invalid("Illegal identifier: %q is a reserved word", string(c.text))
}

func (p *parser) callonFailOnIllegalIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFailOnIllegalIdentifier1()
}

func (c *current) onFailOnMissingDate1() (interface{}, error) {
	return Node{}, invalid("Must include ISO-8601 (YYYY-MM-DD) date as part of timestamp")
}

func (p *parser) callonFailOnMissingDate1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFailOnMissingDate1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// the globalStore allows the parser to store arbitrary values
	globalStore map[string]interface{}
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			globalStore: make(map[string]interface{}),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		p.failAt(true, start.position, ".")
		return p.sliceFrom(start), true
	}
	p.failAt(false, p.pt.position, ".")
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError {
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
