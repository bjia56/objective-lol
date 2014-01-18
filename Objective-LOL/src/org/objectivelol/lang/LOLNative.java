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
			if(toInvoke.getParameterTypes().length != args.length) {
				throw new Exception();
			}
			
			LOLValue result = null;
			
			if(toInvoke.getParameterTypes().length == 0) {
				result = (LOLValue)toInvoke.invoke(this);
			} else {
				result = (LOLValue)toInvoke.invoke(this, (Object[])args);
			}
			
			if(result == null) {
				return LOLNothing.NOTHIN;
			} else {
				return result;
			}
		} catch(Exception e) {
			throw new LOLError("Function with the specified signature not found");
		}
	}
	
}
