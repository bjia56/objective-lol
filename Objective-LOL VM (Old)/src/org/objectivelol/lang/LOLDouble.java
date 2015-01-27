package org.objectivelol.lang;

/**
 * Class to represent a DUBBLE value in Objective-LOL.
 * Essentially a wrapper of a double value.
 * 
 * @author Brett Jia
 */
public class LOLDouble extends LOLNumber {

	/**
	 * Type name as present in Objective-LOL.
	 */
	public static final String TYPE_NAME = "DUBBLE";
	
	private final double value;
	
	/**
	 * Constructor for the LOLDouble class.
	 * 
	 * @param value
	 * A Double representing the value to hold.
	 */
	public LOLDouble(Double value) {
		this.value = value;
	}

	/* (non-Javadoc)
	 * Casts this LOLDouble instance to the specified type.
	 * Can only cast to primitive Objective-LOL types.
	 * 
	 * When casting to INTEGR, any decimal places are truncated.
	 * 
	 * When casting to STRIN, the value is directly converted
	 * into a string of characters.
	 * 
	 * If casting to the specified type is not valid, throws
	 * a LOLError exception.
	 * 
	 * @see org.objectivelol.lang.LOLValue#cast(java.lang.String)
	 */
	@Override
	public LOLValue cast(String type) throws LOLError {
		if(LOLDouble.TYPE_NAME.equals(type) || LOLNumber.TYPE_NAME.equals(type)) {
			return this;
		}
		
		if(LOLInteger.TYPE_NAME.equals(type)) {
			return new LOLInteger((long)value);
		}
		
		if(LOLString.TYPE_NAME.equals(type)) {
			return new LOLString(("" + value).toUpperCase());
		}
		
		throw new LOLError("Cannot cast to the specified type");
	}

	/* (non-Javadoc)
	 * Returns the integer portion of the double value by
	 * truncating any decimal places.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#integerValue()
	 */
	@Override
	public long integerValue() {
		return (long)value;
	}

	/* (non-Javadoc)
	 * Returns the value stored in this instance of LOLDouble.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#doubleValue()
	 */
	@Override
	public double doubleValue() {
		return value;
	}

	/* (non-Javadoc)
	 * Adds the other number to the value stored in this instance of
	 * LOLDouble, preserving any decimal places in both numbers and
	 * upcasting to LOLDouble. The result is returned without changing
	 * either number.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#add(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLNumber add(LOLNumber other) {
		return new LOLDouble(value + other.doubleValue());
	}

	/* (non-Javadoc)
	 * Subtracts the other number from the value stored in this instance
	 * of LOLDouble, preserving any decimal places in both numbers and
	 * upcasting to LOLDouble. The result is returned without changing
	 * either number.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#subtract(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLNumber subtract(LOLNumber other) {
		return new LOLDouble(value - other.doubleValue());
	}

	/* (non-Javadoc)
	 * Multiplies the other number to the value stored in this instance
	 * of LOLDouble, preserving any decimal places in both numbers and
	 * upcasting to LOLDouble. The result is returned without changing
	 * either number.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#multiply(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLNumber multiply(LOLNumber other) {
		return new LOLDouble(value * other.doubleValue());
	}

	/* (non-Javadoc)
	 * Divides the other number from the value stored in this instance
	 * of LOLDouble, preserving any decimal places in both numbers and
	 * upcasting to LOLDouble. The result is returned without changing
	 * either number.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#divide(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLNumber divide(LOLNumber other) {
		return new LOLDouble(value / other.doubleValue());
	}

	/* (non-Javadoc)
	 * Returns the type name as present in Objective-LOL.
	 * 
	 * @see org.objectivelol.lang.LOLValue#getTypeName()
	 */
	@Override
	public String getTypeName() {
		return LOLDouble.TYPE_NAME;
	}

	/* (non-Javadoc)
	 * Checks if the value stored in this instance of LOLDouble
	 * is greater than the other number. Returns a LOLBoolean
	 * containing the boolean value of the comparison.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#greaterThan(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLBoolean greaterThan(LOLNumber other) {
		return (value > other.doubleValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
	}

	/* (non-Javadoc)
	 * Checks if the value stored in this instance of LOLDouble
	 * is less than the other number. Returns a LOLBoolean
	 * containing the boolean value of the comparison.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#lessThan(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLBoolean lessThan(LOLNumber other) {
		return (value < other.doubleValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
	}

	/* (non-Javadoc)
	 * Checks if the value stored in this instance of LOLDouble
	 * is equal to the other number. Returns a LOLBoolean
	 * containing the boolean value of the comparison.
	 * 
	 * @see org.objectivelol.lang.LOLValue#equalTo(org.objectivelol.lang.LOLValue)
	 */
	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		return (value == ((LOLNumber)other.cast(LOLNumber.TYPE_NAME)).doubleValue() ?  LOLBoolean.YEZ : LOLBoolean.NO);
	}

	/* (non-Javadoc)
	 * Duplicates this instance of LOLDouble by constructing
	 * a new LOLDouble object with the double value held
	 * by this.
	 * 
	 * @see org.objectivelol.lang.LOLValue#copy()
	 */
	@Override
	public LOLValue copy() throws LOLError {
		return new LOLDouble(value);
	}

}
