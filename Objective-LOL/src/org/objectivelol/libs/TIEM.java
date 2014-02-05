package org.objectivelol.libs;

import org.objectivelol.lang.LOLInteger;
import org.objectivelol.lang.LOLNative;
import org.objectivelol.lang.LOLValue;

public class TIEM extends LOLNative {

	public static LOLInteger NAO() {
		return (LOLInteger)LOLValue.valueOf(System.currentTimeMillis());
	}
	
}
