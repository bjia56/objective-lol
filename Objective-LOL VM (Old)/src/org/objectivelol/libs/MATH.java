package org.objectivelol.libs;

import org.objectivelol.lang.LOLDouble;
import org.objectivelol.lang.LOLInteger;
import org.objectivelol.lang.LOLNative;
import org.objectivelol.lang.LOLNumber;
import org.objectivelol.lang.LOLValue;

public class MATH extends LOLNative {

	public static LOLDouble ACOS(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.acos(arg.doubleValue())));
	}

	public static LOLDouble ASIN(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.asin(arg.doubleValue())));
	}

	public static LOLDouble ATAN(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.atan(arg.doubleValue())));
	}

	public static LOLDouble ATAN2(LOLNumber arg1, LOLNumber arg2) {
		return (LOLDouble)LOLValue.valueOf((Math.atan2(arg1.doubleValue(), arg2.doubleValue())));
	}

	public static LOLInteger BITAN(LOLInteger arg1, LOLInteger arg2) {
		return (LOLInteger)LOLValue.valueOf(arg1.integerValue() & arg2.integerValue());
	}
	
	public static LOLInteger BITOR(LOLInteger arg1, LOLInteger arg2) {
		return (LOLInteger)LOLValue.valueOf(arg1.integerValue() | arg2.integerValue());
	}
	
	public static LOLInteger BITXOR(LOLInteger arg1, LOLInteger arg2) {
		return (LOLInteger)LOLValue.valueOf(arg1.integerValue() ^ arg2.integerValue());
	}
	
	public static LOLDouble CBRT(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.cbrt(arg.doubleValue())));
	}

	public static LOLDouble COS(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.cos(arg.doubleValue())));
	}

	public static LOLDouble COSH(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.cosh(arg.doubleValue())));
	}

	public static LOLDouble EXP(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.exp(arg.doubleValue())));
	}

	public static LOLDouble LOG(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.log(arg.doubleValue())));
	}

	public static LOLDouble LOG10(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.log10(arg.doubleValue())));
	}

	public static LOLDouble POW(LOLNumber arg1, LOLNumber arg2) {
		return (LOLDouble)LOLValue.valueOf((Math.pow(arg1.doubleValue(), arg2.doubleValue())));
	}

	public static LOLDouble RAND() {
		return (LOLDouble)LOLValue.valueOf(Math.random());
	}

	public static LOLDouble SIN(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.sin(arg.doubleValue())));
	}

	public static LOLDouble SINH(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.sinh(arg.doubleValue())));
	}

	public static LOLDouble SQRT(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.sqrt(arg.doubleValue())));
	}

	public static LOLDouble TAN(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.tan(arg.doubleValue())));
	}

	public static LOLDouble TANH(LOLNumber arg) {
		return (LOLDouble)LOLValue.valueOf((Math.tanh(arg.doubleValue())));
	}

}
