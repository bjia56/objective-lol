package org.objectivelol.lang;

/**
 * Abstract class to represent any NUMBR value in
 * Objective-LOL. All numeric classes extend this.
 * 
 * @author Brett Jia
 */
public abstract class LOLNumber extends LOLValue {

	/**
	 * Type name as present in Objective-LOL.
	 */
	public static final String TYPE_NAME = "NUMBR";

	/**
	 * Returns the value stored in this LOLNumber
	 * as a Java integer (long) value. Casts or
	 * truncations are performed as needed.
	 * 
	 * @return
	 * A long representing the integer portion of
	 * the value stored.
	 */
	public abstract long integerValue();

	/**
	 * Returns the value stored in this LOLNumber
	 * as a Java double value. Casts are performed
	 * as needed.
	 * 
	 * @return
	 * A double representing the value stored.
	 */
	public abstract double doubleValue();

	/**
	 * Adds the value stored in other to the value
	 * stored in this LOLNumber, then returns the
	 * result (this + other). Neither LOLNumber value
	 * is changed.
	 * 
	 * The LOLNumber returned is of the highest
	 * precision necessary (i.e. a LOLInteger and a
	 * LOLDouble would return a LOLDouble, while
	 * a LOLInteger and a LOLInteger would only return
	 * a LOLInteger).
	 * 
	 * @param other
	 * A LOLNumber representing the other number to
	 * add.
	 * 
	 * @return
	 * A LOLNumber representing the result of the
	 * addition operation.
	 */
	public abstract LOLNumber add(LOLNumber other);

	/**
	 * Subtracts the value stored in other from the
	 * value stored in this LOLNumber, then returns
	 * the result (this - other). Neither LOLNumber
	 * value is changed.
	 * 
	 * The LOLNumber returned is of the highest
	 * precision necessary (i.e. a LOLInteger and a
	 * LOLDouble would return a LOLDouble, while
	 * a LOLInteger and a LOLInteger would only return
	 * a LOLInteger).
	 * 
	 * @param other
	 * A LOLNumber representing the other number to
	 * subtract.
	 * 
	 * @return
	 * A LOLNumber representing the result of the
	 * subtraction operation.
	 */
	public abstract LOLNumber subtract(LOLNumber other);

	/**
	 * Multiplies the value stored in other with the
	 * value stored in this LOLNumber, then returns
	 * the result (this * other). Neither LOLNumber
	 * value is changed.
	 * 
	 * The LOLNumber returned is of the highest
	 * precision necessary (i.e. a LOLInteger and a
	 * LOLDouble would return a LOLDouble, while
	 * a LOLInteger and a LOLInteger would only return
	 * a LOLInteger).
	 * 
	 * @param other
	 * A LOLNumber representing the other number to
	 * multiply.
	 * 
	 * @return
	 * A LOLNumber representing the result of the
	 * multiplication operation.
	 */
	public abstract LOLNumber multiply(LOLNumber other);

	/**
	 * Divides the value stored in other from the
	 * value stored in this LOLNumber, then returns
	 * the result (this / other). Neither LOLNumber
	 * value is changed.
	 * 
	 * The LOLNumber returned is of the highest
	 * precision necessary (i.e. a LOLInteger and a
	 * LOLDouble would return a LOLDouble, while
	 * a LOLInteger and a LOLInteger would only return
	 * a LOLInteger).
	 * 
	 * @param other
	 * A LOLNumber representing the other number to
	 * divide.
	 * 
	 * @return
	 * A LOLNumber representing the result of the
	 * division operation.
	 */
	public abstract LOLNumber divide(LOLNumber other);

	/**
	 * Checks if the value stored in this LOLNumber
	 * is greater than the value stored in other.
	 * Equivalent to the test (this > other).
	 * 
	 * @param other
	 * A LOLNumber representing the other number to
	 * compare.
	 * 
	 * @return
	 * A LOLBoolean representing the result of the
	 * comparison.
	 */
	public abstract LOLBoolean greaterThan(LOLNumber other);

	/**
	 * Checks if the value stored in this LOLNumber
	 * is less than the value stored in other.
	 * Equivalent to the test (this < other).
	 * 
	 * @param other
	 * A LOLNumber representing the other number to
	 * compare.
	 * 
	 * @return
	 * A LOLBoolean representing the result of the
	 * comparison.
	 */
	public abstract LOLBoolean lessThan(LOLNumber other);

}
