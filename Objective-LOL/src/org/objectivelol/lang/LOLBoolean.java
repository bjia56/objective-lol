package org.objectivelol.lang;

/**
 * Class to represent a BOOL value in Objective-LOL.
 * Essentially a wrapper of a boolean value.
 * 
 * @author Brett Jia
 */
public class LOLBoolean extends LOLValue {

	/**
	 * Type name as present in Objective-LOL.
	 */
	public static final String TYPE_NAME = "BOOL";
	
	/**
	 * LOLBoolean instance representative of the boolean
	 * value of true, or YEZ in Objective-LOL
	 */
	public static final LOLBoolean YEZ = new LOLBoolean(true);
	
	/**
	 * LOLBoolean instance representative of the boolean
	 * value of false, or NO in Objective-LOL
	 */
	public static final LOLBoolean NO = new LOLBoolean(false);
	
	private final boolean value;
	
	/**
	 * Constructor for the LOLBoolean class.
	 * 
	 * @param value
	 * A Boolean representing the value to hold.
	 */
	private LOLBoolean(Boolean value) {
		this.value = value;
	}

	/* (non-Javadoc)
	 * Casts this LOLBoolean instance to the specified type.
	 * Can only cast to primitive Objective-LOL types.
	 * 
	 * When casting to INTEGR or DUBBLE, YEZ becomes
	 * 1 and NO becomes 0.
	 * 
	 * When casting to STRIN, YEZ becomes "YEZ" and
	 * NO becomes "NO".
	 * 
	 * If casting to the specified type is not valid, throws
	 * a LOLError exception.
	 * 
	 * @see org.objectivelol.lang.LOLValue#cast(java.lang.String)
	 */
	@Override
	public LOLValue cast(String type) throws LOLError {
		if(LOLInteger.TYPE_NAME.equals(type) || LOLNumber.TYPE_NAME.equals(type)) {
			if(value) {
				return new LOLInteger(1l);
			} else {
				return new LOLInteger(0l);
			}
		}
		
		if(LOLDouble.TYPE_NAME.equals(type)) {
			if(value) {
				return new LOLDouble(1d);
			} else {
				return new LOLDouble(0d);
			}
		}
		
		if(LOLString.TYPE_NAME.equals(type)) {
			return new LOLString((value ? "YEZ" : "NO").toUpperCase());
		}
		
		if(LOLBoolean.TYPE_NAME.equals(type)) {
			return this;
		}
		
		throw new LOLError("Cannot cast to the specified type");
	}
	
	/**
	 * Gives the Java boolean value held by this instance of
	 * LOLBoolean.
	 * 
	 * @return
	 * The boolean value held.
	 */
	public boolean booleanValue() {
		return value;
	}

	/* (non-Javadoc)
	 * Returns the type name as present in Objective-LOL.
	 * 
	 * @see org.objectivelol.lang.LOLValue#getTypeName()
	 */
	@Override
	public String getTypeName() {
		return LOLBoolean.TYPE_NAME;
	}

	/* (non-Javadoc)
	 * Returns whether the boolean value held by this
	 * instance of LOLBoolean is equal to the value held
	 * by the result of casting the other LOLValue to
	 * BOOL.
	 * 
	 * @see org.objectivelol.lang.LOLValue#equalTo(org.objectivelol.lang.LOLValue)
	 */
	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		return (value == ((LOLBoolean)other.cast(LOLBoolean.TYPE_NAME)).value ? LOLBoolean.YEZ : LOLBoolean.NO);
	}

	/* (non-Javadoc)
	 * Duplicates this instance of LOLBoolean by constructing
	 * a new LOLBoolean object with the boolean value held
	 * by this.
	 * 
	 * @see org.objectivelol.lang.LOLValue#copy()
	 */
	@Override
	public LOLValue copy() throws LOLError {
		return new LOLBoolean(value);
	}

}
