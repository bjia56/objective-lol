package org.objectivelol.lang;

public class LOLInteger extends LOLNumber {

	public static final String TYPE_NAME = "INTEGR";
	
	private final long value;
	
	public LOLInteger(Long value) {
		this.value = value;
	}
	
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

	@Override
	public long integerValue() {
		return value;
	}

	@Override
	public double doubleValue() {
		return value;
	}

	@Override
	public LOLNumber add(LOLNumber other) {
		if(other.isLOLDouble()) {
			return new LOLDouble(value + other.doubleValue());
		} else {
			return new LOLInteger(value + other.integerValue());
		}
	}

	@Override
	public LOLNumber subtract(LOLNumber other) {
		if(other.isLOLDouble()) {
			return new LOLDouble(value - other.doubleValue());
		} else {
			return new LOLInteger(value - other.integerValue());
		}
	}

	@Override
	public LOLNumber multiply(LOLNumber other) {
		if(other.isLOLDouble()) {
			return new LOLDouble(value * other.doubleValue());
		} else {
			return new LOLInteger(value * other.integerValue());
		}
	}

	@Override
	public LOLNumber divide(LOLNumber other) {
		if(other.isLOLDouble()) {
			return new LOLDouble(value / other.doubleValue());
		} else {
			return new LOLInteger(value / other.integerValue());
		}
	}

	@Override
	public String getTypeName() {
		return LOLInteger.TYPE_NAME;
	}

	@Override
	public LOLBoolean greaterThan(LOLNumber other) {
		if(other.isLOLDouble()) {
			return(value > other.doubleValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		} else {
			return (value > other.integerValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		}
	}

	@Override
	public LOLBoolean lessThan(LOLNumber other) {
		if(other.isLOLDouble()) {
			return (value < other.doubleValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		} else {
			return (value < other.integerValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		}
	}

	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		LOLNumber ln = (LOLNumber)other.cast(LOLNumber.TYPE_NAME);
		
		if(ln.isLOLDouble()) {
			return(value == ln.doubleValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		} else {
			return (value == ln.integerValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
		}
	}
	
}
