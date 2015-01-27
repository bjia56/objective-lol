package org.objectivelol.lang;

/**
 * Class to represent the NOTHIN value in Objective-LOL.
 * 
 * @author Brett Jia
 */
public class LOLNothing extends LOLValue {
	
	/**
	 * Type name as present in Objective-LOL.
	 */
	public static final String TYPE_NAME = "NOTHIN";
	
	/**
	 * Global instance of LOLNothing.
	 */
	public static final LOLNothing NOTHIN = new LOLNothing();

	/* (non-Javadoc)
	 * Cannot cast LOLNothing to any other type, so throws a
	 * LOLError exception.
	 * 
	 * @see org.objectivelol.lang.LOLValue#cast(java.lang.String)
	 */
	@Override
	public LOLValue cast(String type) throws LOLError {
		throw new LOLError("Cannot cast to the specified type");
	}

	/*
	 * (non-Javadoc)
	 * Returns the type name as present in Objective-LOL.
	 * 
	 * @see org.objectivelol.lang.LOLValue#getTypeName()
	 */
	@Override
	public String getTypeName() {
		return LOLNothing.TYPE_NAME;
	}

	/*
	 * (non-Javadoc)
	 * Checks if the value stored in other is equal to
	 * LOLNothing. Only returns the LOLBoolean equivalent
	 * of true if other is a LOLNothing, false otherwise.
	 * 
	 * @see org.objectivelol.lang.LOLValue#equalTo(org.objectivelol.lang.LOLValue)
	 */
	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		return (other.isLOLNothing() ? LOLBoolean.YEZ : LOLBoolean.NO);
	}

	/*
	 * (non-Javadoc)
	 * Function used when attempting to duplicate a LOLValue.
	 * Since a LOLNothing is (more or less) a constant value,
	 * this function returns a reference to this.
	 * 
	 * @see org.objectivelol.lang.LOLValue#copy()
	 */
	@Override
	public LOLValue copy() throws LOLError {
		return this;
	}
	
}
