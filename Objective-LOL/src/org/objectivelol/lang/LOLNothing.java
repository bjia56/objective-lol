package org.objectivelol.lang;

public class LOLNothing extends LOLValue {
	
	public static final String TYPE_NAME = "NOTHIN";
	
	public static final LOLNothing NOTHIN = new LOLNothing();

	@Override
	public LOLValue cast(String type) throws LOLError {
		throw new LOLError("Cannot cast to the specified type");
	}

	@Override
	public String getTypeName() {
		return LOLNothing.TYPE_NAME;
	}

	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		return (other.isLOLNothing() ? LOLBoolean.YEZ : LOLBoolean.NO);
	}
	
}
