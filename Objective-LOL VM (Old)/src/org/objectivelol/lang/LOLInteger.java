package org.objectivelol.lang;

/**
 * Class to represent an INTEGR value in Objective-LOL.
 * Essentially a wrapper of a long value.
 * 
 * @author Brett Jia
 */
public class LOLInteger extends LOLNumber {

	/**
	 * Type name as present in Objective-LOL.
	 */
	public static final String TYPE_NAME = "INTEGR";
	
	private final long value;
	
	/**
	 * Constructor for the LOLInteger class.
	 * 
	 * @param value
	 * A Long representing the value to hold.
	 */
	public LOLInteger(Long value) {
		this.value = value;
	}
	
	/* (non-Javadoc)
	 * Casts this LOLInteger instance to the specified type.
	 * Can only cast to primitive Objective-LOL types.
	 * 
	 * When casting to DUBBLE, the value is upcasted.
	 * 
	 * When casting to STRIN, the value is directly converted
	 * into a string of characters.
	 * 
	 * When casting to BOOL, the result is NO if the value is 0,
	 * or YEZ otherwise.
	 * 
	 * If casting to the specified type is not valid, throws
	 * a LOLError exception.
	 * 
	 * @see org.objectivelol.lang.LOLValue#cast(java.lang.String)
	 */
	@Override
	public LOLValue cast(String type) throws LOLError {
		if(LOLInteger.TYPE_NAME.equals(type) || LOLNumber.TYPE_NAME.equals(type)) {
			return this;
		}
		
		if(LOLDouble.TYPE_NAME.equals(type)) {
			return new LOLDouble((double)value);
		}
		
		if(LOLString.TYPE_NAME.equals(type)) {
			return new LOLString(("" + value).toUpperCase());
		}
		
		if(LOLBoolean.TYPE_NAME.equals(type)) {
			if(value == 0) {
				return LOLBoolean.NO;
			} else {
				return LOLBoolean.YEZ;
			}
		}
		
		throw new LOLError("Cannot cast to the specified type");
	}

	/* (non-Javadoc)
	 * Returns the value stored in this instance of LOLInteger.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#integerValue()
	 */
	@Override
	public long integerValue() {
		return value;
	}

	/* (non-Javadoc)
	 * Upcasts the value stored and returns the result as a
	 * double.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#doubleValue()
	 */
	@Override
	public double doubleValue() {
		return value;
	}

	/* (non-Javadoc)
	 * Adds the other number to the value stored in this instance of
	 * LOLInteger, preserving any decimal places in the other number
	 * if necessary and returning a LOLNumber denoting the sum. The
	 * result is returned without changing either number.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#add(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLNumber add(LOLNumber other) {
		if(other.isLOLDouble()) {
			return new LOLDouble(value + other.doubleValue());
		} else {
			return new LOLInteger(value + other.integerValue());
		}
	}

	/* (non-Javadoc)
	 * Subtracts the other number from the value stored in this instance
	 * of LOLInteger, preserving any decimal places in the other number
	 * if necessary and returning a LOLNumber denoting the difference.
	 * The result is returned without changing either number.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#add(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLNumber subtract(LOLNumber other) {
		if(other.isLOLDouble()) {
			return new LOLDouble(value - other.doubleValue());
		} else {
			return new LOLInteger(value - other.integerValue());
		}
	}

	/* (non-Javadoc)
	 * Multiplies the other number to the value stored in this instance
	 * of LOLInteger, preserving any decimal places in the other number
	 * if necessary and returning a LOLNumber denoting the product.
	 * The result is returned without changing either number.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#add(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLNumber multiply(LOLNumber other) {
		if(other.isLOLDouble()) {
			return new LOLDouble(value * other.doubleValue());
		} else {
			return new LOLInteger(value * other.integerValue());
		}
	}

	/* (non-Javadoc)
	 * Divides the other number from the value stored in this instance
	 * of LOLInteger, preserving any decimal places in the other number
	 * if necessary and returning a LOLNumber denoting the quotient.
	 * The result is returned without changing either number.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#add(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLNumber divide(LOLNumber other) {
		if(other.isLOLDouble()) {
			return new LOLDouble(value / other.doubleValue());
		} else {
			return new LOLInteger(value / other.integerValue());
		}
	}

	/* (non-Javadoc)
	 * Returns the type name as present in Objective-LOL.
	 * 
	 * @see org.objectivelol.lang.LOLValue#getTypeName()
	 */
	@Override
	public String getTypeName() {
		return LOLInteger.TYPE_NAME;
	}

	/* (non-Javadoc)
	 * Checks if the value stored in this instance of LOLInteger
	 * is greater than the other number. Returns a LOLBoolean
	 * containing the boolean value of the comparison.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#greaterThan(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLBoolean greaterThan(LOLNumber other) {
		if(other.isLOLDouble()) {
			return(value > other.doubleValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		} else {
			return (value > other.integerValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		}
	}

	/* (non-Javadoc)
	 * Checks if the value stored in this instance of LOLInteger
	 * is less than the other number. Returns a LOLBoolean
	 * containing the boolean value of the comparison.
	 * 
	 * @see org.objectivelol.lang.LOLNumber#lessThan(org.objectivelol.lang.LOLNumber)
	 */
	@Override
	public LOLBoolean lessThan(LOLNumber other) {
		if(other.isLOLDouble()) {
			return (value < other.doubleValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		} else {
			return (value < other.integerValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		}
	}

	/* (non-Javadoc)
	 * Checks if the value stored in this instance of LOLInteger
	 * is equal to the other number. Returns a LOLBoolean
	 * containing the boolean value of the comparison.
	 * 
	 * @see org.objectivelol.lang.LOLValue#equalTo(org.objectivelol.lang.LOLValue)
	 */
	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		LOLNumber ln = (LOLNumber)other.cast(LOLNumber.TYPE_NAME);
		
		if(ln.isLOLDouble()) {
			return(value == ln.doubleValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		} else {
			return (value == ln.integerValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		}
	}

	/* (non-Javadoc)
	 * Duplicates this instance of LOLInteger by constructing
	 * a new LOLInteger object with the long value held
	 * by this.
	 * 
	 * @see org.objectivelol.lang.LOLValue#copy()
	 */
	@Override
	public LOLValue copy() throws LOLError {
		return new LOLInteger(value);
	}
	
}
