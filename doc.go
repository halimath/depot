// Package depot implements a small abstraction layer for accessing sql (and
// potentially other) databases.
//
// depot is build around two concepts: values and clauses. See the respective
// descriptions for details.
//
// depot also provides a code generator used to generate repository types
// from regulare go structs with field tags. The resulting code uses no
// reflection and provides a typesafe interface to interact with the database.
package depot
