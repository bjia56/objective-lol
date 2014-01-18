package org.objectivelol.lang;

public class LOLString extends LOLValue {

	public static final String TYPE_NAME = "STRIN";

	private final String value;

	public LOLString(String value) {
		this.value = value;
	}

	@Override
	public LOLValue cast(String type) throws LOLError {
		if(LOLDouble.TYPE_NAME.equals(type) || LOLNumber.TYPE_NAME.equals(type)) {
			try {
				return new LOLDouble(Double.parseDouble(value));
			} catch(ArithmeticException e) {
				throw new LOLError("Cannot cast to the specified type");
			}
		}

		if(LOLInteger.TYPE_NAME.equals(type)) {
			try {
				return new LOLInteger(Long.parseLong(value));
			} catch(ArithmeticException e) {
				try {
					return new LOLInteger(Long.parseLong(value.toUpperCase().replaceFirst("0X", ""), 16));
				} catch(ArithmeticException e2) {
					throw new LOLError("Cannot cast to the specified type");
				}
			}
		}

		if(LOLBoolean.TYPE_NAME.equals(type)) {
			if("YEZ".equals(value)) {
				return LOLBoolean.YEZ;
			} else {
				return LOLBoolean.NO;
			}
		}

		if(LOLString.TYPE_NAME.equals(type)) {
			return this;
		}

		throw new LOLError("Cannot cast to the specified type");
	}

	@Override
	public String getTypeName() {
		return LOLString.TYPE_NAME;
	}

	@Override
	public String toString() {
		return value;
	}

	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		return (value.equals(((LOLString)other.cast(LOLString.TYPE_NAME)).toString()) ?  LOLBoolean.YEZ : LOLBoolean.NO);
	}

}
