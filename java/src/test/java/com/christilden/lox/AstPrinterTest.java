package com.christilden.lox;

import junit.framework.Test;
import junit.framework.TestCase;
import junit.framework.TestSuite;

/**
 * Unit test for AstPrinter.
 */
public class AstPrinterTest extends TestCase {
    /**
     * Create the test case
     *
     * @param testName name of the test case
     */
    public AstPrinterTest( String testName ) {
        super( testName );
    }

    /**
     * @return the suite of tests being tested
     */
    public static Test suite() {
        return new TestSuite( AstPrinterTest.class );
    }

    /**
     * Tests the AST printer.
     */
    public void testAstPrinter() {
        Expr expression = new Expr.Binary(
            new Expr.Unary(
                new Token(TokenType.MINUS, "-", null, 1),
                new Expr.Literal(123)),
            new Token(TokenType.STAR, "*", null, 1),
            new Expr.Grouping(
                new Expr.Literal(45.67)));

        assertEquals(new AstPrinter().print(expression), "(* (- 123) (group 45.67))");
    }
}
