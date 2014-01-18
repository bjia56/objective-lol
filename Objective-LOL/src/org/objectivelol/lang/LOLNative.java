package org.objectivelol.lang;

import java.lang.reflect.Method;

public abstract class LOLNative {

	public final LOLValue invoke(String methodName, LOLValue ... args) throws LOLError {
		Method toInvoke = null;
		
		for(Method m : this.getClass().getMethods()) {
			if(m.getName().equals(methodName)) {
				toInvoke = m;
				break;
			}
		}
		
		try {
			return (LOLValue)toInvoke.invoke(this, (Object[])args);
		} catch (Exception e) {
			throw new LOLError("Function with the specified signature not found");
		}
	}
	
}
