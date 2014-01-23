package org.objectivelol.lang;

public class LOLDouble extends LOLNumber {

	public static final String TYPE_NAME = "DUBBLE";
	
	private final double value;
	
	public LOLDouble(Double value) {
		this.value = value;
	}

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

	@Override
	public long integerValue() {
		return (long)value;
	}

	@Override
	public double doubleValue() {
		return value;
	}

	@Override
	public LOLNumber add(LOLNumber other) {
		return new LOLDouble(value + other.doubleValue());
	}

	@Override
	public LOLNumber subtract(LOLNumber other) {
		return new LOLDouble(value - other.doubleValue());
	}

	@Override
	public LOLNumber multiply(LOLNumber other) {
		return new LOLDouble(value * other.doubleValue());
	}

	@Override
	public LOLNumber divide(LOLNumber other) {
		return new LOLDouble(value / other.doubleValue());
	}

	@Override
	public String getTypeName() {
		return LOLDouble.TYPE_NAME;
	}

	@Override
	public LOLBoolean greaterThan(LOLNumber other) {
		return (value > other.doubleValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
	}

	@Override
	public LOLBoolean lessThan(LOLNumber other) {
		return (value < other.doubleValue() ? LOLBoolean.YEZ : LOLBoolean.NO);
	}

	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		return (value == ((LOLNumber)other.cast(LOLNumber.TYPE_NAME)).doubleValue() ?  LOLBoolean.YEZ : LOLBoolean.NO);
	}

	@Override
	public LOLValue copy() throws LOLError {
		return new LOLDouble(value);
	}

}
