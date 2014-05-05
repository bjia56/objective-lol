package org.objectivelol.lang;

import java.lang.reflect.Method;

/**
 * Abstract class to represent a class of NATIV functions in
 * Objective-LOL. Implements the Java reflection utilities
 * necessary to look up and run the native function based on
 * name. 
 * 
 * @author Brett Jia
 */
public abstract class LOLNative {

	/**
	 * Invokes the Java function defined in a subclass of LOLNative
	 * by looking it up by name. The function must not be overloaded.
	 * 
	 * @param methodName
	 * A String representing the name of the method to invoke.
	 * 
	 * @param args
	 * A LOLValue varargs specifying the arguments to be passed into
	 * the function call. These arguments must match in number to the
	 * arguments of the function to invoke, as well as sequence of
	 * types.
	 * 
	 * @return
	 * Returns the result of running the method, if successful. If
	 * no return is given by the function, this function returns
	 * LOLNothing.
	 * 
	 * @throws LOLError
	 * Throws a LOLError if the specified method is not found, if
	 * there is an argument mismatch, or if there is an error in
	 * execution of the specified Java function.
	 */
	public final LOLValue invoke(String methodName, LOLValue ... args) throws LOLError {
		Method toInvoke = null;
		
		// find the method name, if it exists
		for(Method m : this.getClass().getMethods()) {
			if(m.getName().equals(methodName)) {
				toInvoke = m;
				break;
			}
		}
		
		try {
			// check if the number of parameters match
			if(toInvoke.getParameterTypes().length != args.length) {
				throw new Exception();
			}
			
			LOLValue result = null;
			
			// try to invoke the method
			if(toInvoke.getParameterTypes().length == 0) {
				result = (LOLValue)toInvoke.invoke(this);
			} else {
				result = (LOLValue)toInvoke.invoke(this, (Object[])args);
			}
			
			// check if there was a return value, and return accordingly
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
