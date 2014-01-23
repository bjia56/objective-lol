package org.objectivelol.lang;

public class LOLBoolean extends LOLValue {

	public static final String TYPE_NAME = "BOOL";
	
	public static final LOLBoolean YEZ = new LOLBoolean(true);
	public static final LOLBoolean NO = new LOLBoolean(false);
	
	private final boolean value;
	
	private LOLBoolean(Boolean value) {
		this.value = value;
	}

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
	
	public boolean booleanValue() {
		return value;
	}

	@Override
	public String getTypeName() {
		return LOLBoolean.TYPE_NAME;
	}

	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		return (value == ((LOLBoolean)other.cast(LOLBoolean.TYPE_NAME)).value ? LOLBoolean.YEZ : LOLBoolean.NO);
	}

	@Override
	public LOLValue copy() throws LOLError {
		return new LOLBoolean(value);
	}

}
