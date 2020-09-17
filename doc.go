/*

A typical use of the errxpect package is by importing it into the in the current
file's file block for easy reference without needing the package name.

Before, Gomega forced you to write assertions where you are only interested in
checking the error return value of a multi-return value function:

    // func multifoo() (string, int, error) { ... }

    _, _, err := multifoo()
    Expect(err).To(Succeed())

This can now be rewritten in a more concise form as:

    Errxpect(multifoo()).To(Succeed())

ExpectWithError() needs to be slightly phrased differently, due to Golang's
language restrictions, but probably even more neat now:

    Errpect(multifoo()).WithOffset(1).To(Succeed())

*/
package errxpect
