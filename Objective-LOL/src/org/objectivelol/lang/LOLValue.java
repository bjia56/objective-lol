package org.objectivelol.lang;

public abstract class LOLValue {
	
	public static class ValueStruct {
		
		private final String type;
		private LOLValue value;
		private final boolean isLocked;
		
		public ValueStruct(String type, LOLValue value, boolean isLocked) {
			this.type = type;
			this.value = value;
			this.isLocked = isLocked;
		}
		
		public ValueStruct copy() {
			return new ValueStruct(type, value, isLocked);
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

	public static LOLValue valueOf(Object o) {
		if(o instanceof LOLValue) {
			return (LOLValue)o;
		}

		if(o instanceof Integer || o instanceof Long) {
			return new LOLInteger(((Number)o).longValue());
		}

		if(o instanceof Double || o instanceof Float) {
			return new LOLDouble(((Number)o).doubleValue());
		}

		if(o instanceof Boolean) {
			return ((Boolean)o ? LOLBoolean.YEZ : LOLBoolean.NO);
		}
		
		if(o instanceof Character) {
			return new LOLString((Character)o + "");
		}

		if(o instanceof String) {
			String str = (String)o;
			
			try {
				return new LOLInteger(Long.parseLong(str));
			} catch(NumberFormatException e) {
				try {
					String str2 = str.toUpperCase();
					if(!str2.startsWith("0X")) {
						throw new NumberFormatException();
					}
					return new LOLInteger(Long.parseLong(str2.replaceFirst("0X", ""), 16));
				} catch(NumberFormatException e2) {
					try {
						return new LOLDouble(Double.parseDouble(str));
					} catch(NumberFormatException e3) {
						if(str.equals("YEZ")) {
							return LOLBoolean.YEZ;
						}
						if(str.equals("NO")) {
							return LOLBoolean.NO;
						}
					}
				}
			}
			
			return new LOLString(str);
		}

		throw new IllegalArgumentException("Argument cannot be converted to a primitive Objective-LOL type");
	}
	
	public boolean isLOLNothing() {
		return this instanceof LOLNothing;
	}
	
	public boolean isLOLDouble() {
		return this instanceof LOLDouble;
	}

	public boolean isLOLInteger() {
		return this instanceof LOLInteger;
	}
	
	public boolean isLOLBoolean() {
		return this instanceof LOLBoolean;
	}
	
	public boolean isLOLNumber() {
		return this instanceof LOLNumber;
	}
	
	public boolean isLOLString() {
		return this instanceof LOLString;
	}
	
	public abstract LOLValue cast(String type) throws LOLError;
	
	public abstract String getTypeName();

	public abstract LOLBoolean equalTo(LOLValue other) throws LOLError;
	
}
