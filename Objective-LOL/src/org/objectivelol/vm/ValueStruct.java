package org.objectivelol.vm;

import org.objectivelol.lang.LOLError;
import org.objectivelol.lang.LOLValue;

public class ValueStruct {
	
	private final String type;
	private LOLValue value;
	private final boolean isLocked;
	
	public ValueStruct(String type, LOLValue value, boolean isLocked) {
		this.type = type;
		this.value = value;
		this.isLocked = isLocked;
	}
	
	public ValueStruct copy() throws LOLError {
		return new ValueStruct(type, value.copy(), isLocked);
	}
	
	public String getType() {
		return type;
	}
	
	public LOLValue getValue() {
		return value;
	}
	
	public boolean getIsLocked() {
		return isLocked;
	}
	
	public void setValue(LOLValue newValue) throws LOLError {
		if(isLocked) {
			throw new LOLError("Variable is locked and cannot be updated");
		}
		
		value = newValue.cast(type);
	}
	
}
