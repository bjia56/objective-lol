package org.objectivelol.lang;

public abstract class LOLValue {

	/**
	 * Converts an arbitrary Java object into a LOLValue.
	 * Conversion is currently limited to Java primitives.
	 * Any unsupported conversions will throw an exception.
	 * 
	 * @param o
	 * An Object representing the value to convert.
	 * 
	 * @return
	 * A LOLValue representing the result of the conversion.
	 * The conversion pattern is as follows:
	 * <li>LOLValues will be returned directly</li>
	 * <li>Integers and Longs will be converted into
	 * LOLIntegers</li>
	 * <li>Doubles and Floats will be converted into
	 * LOLDoubles</li>
	 * <li>Booleans will be converted into LOLBooleans</li>
	 * <li>Characters will be converted into LOLStrings of
	 * one character long</li>
	 * <li>Strings will be converted into a LOLInteger,
	 * LOLDouble, or LOLBoolean, if possible; otherwise, a LOLString will be returned</li> 
	 */
	public static LOLValue valueOf(Object o) {
		if(o instanceof LOLValue) {
			// no reason to do anything special here, just return
			// the value
			return (LOLValue)o;
		}

		if(o instanceof Integer || o instanceof Long) {
			// convert Integers and Longs into LOLIntegers
			return new LOLInteger(((Number)o).longValue());
		}

		if(o instanceof Double || o instanceof Float) {
			// convert Doubles and Floats into LOLDoubles
			return new LOLDouble(((Number)o).doubleValue());
		}

		if(o instanceof Boolean) {
			// converts Booleans into the LOLBoolean YEZ and
			// NO constants
			return ((Boolean)o ? LOLBoolean.YEZ : LOLBoolean.NO);
		}
		
		if(o instanceof Character) {
			// converts Characters into a one character String
			return new LOLString((Character)o + "");
		}

		if(o instanceof String) {
			String str = (String)o;
			
			// Strings can be converted into a wide variety of LOLValues,
			// depending on the contents of the String
			try {
				// try to convert the String into a Long
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
	
	public abstract LOLValue copy() throws LOLError;
	
}
