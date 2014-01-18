package org.objectivelol.lang;

public abstract class LOLNumber extends LOLValue {
	
	public static final String TYPE_NAME = "NUMBR";
	
	public abstract long integerValue();
	
	public abstract double doubleValue();
	
	public abstract LOLNumber add(LOLNumber other);
	
	public abstract LOLNumber subtract(LOLNumber other);
	
	public abstract LOLNumber multiply(LOLNumber other);
	
	public abstract LOLNumber divide(LOLNumber other);
	
	public abstract LOLBoolean greaterThan(LOLNumber other);
	
	public abstract LOLBoolean lessThan(LOLNumber other);
	
}
