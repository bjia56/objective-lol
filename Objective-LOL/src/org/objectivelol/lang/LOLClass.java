package org.objectivelol.lang;

import java.util.Collection;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map.Entry;

import org.objectivelol.vm.ValueStruct;

/**
 * Class to represent CLAS configurations in
 * Objective-LOL. Used to instantiate new
 * objects during runtime.
 * 
 * @author Brett Jia
 */
public class LOLClass {

	private final String className;

	private final HashMap<String, ValueStruct> publicMemberVariables;
	private final HashMap<String, ValueStruct> privateMemberVariables;
	private final HashMap<String, ValueStruct> publicSharedVariables;
	private final HashMap<String, ValueStruct> privateSharedVariables;

	private final HashMap<String, LOLFunction> publicMemberFunctions;
	private final HashMap<String, LOLFunction> privateMemberFunctions;
	private final HashMap<String, LOLFunction> publicSharedFunctions;
	private final HashMap<String, LOLFunction> privateSharedFunctions;

	private final String parentSource;

	/**
	 * Constructor for the LOLClass class.
	 * 
	 * @param className
	 * A String representing the name of the CLAS.
	 * 
	 * @param publicMemberVariables
	 * A HashMap of String to ValueStruct representing
	 * the public member variables of the CLAS.
	 * 
	 * @param privateMemberVariables
	 * A HashMap of String to ValueStruct representing
	 * the private member variables of the CLAS.
	 * 
	 * @param publicSharedVariables
	 * A HashMap of String to ValueStruct representing
	 * the public SHARD variables of the CLAS.
	 * 
	 * @param privateSharedVariables
	 * A HashMap of String to ValueStruct representing
	 * the private SHARD variables of the CLAS.
	 * 
	 * @param publicMemberFunctions
	 * A HashMap of String to LOLFunction representing
	 * the public member functions of the CLAS.
	 * 
	 * @param privateMemberFunctions
	 * A HashMap of String to LOLFunction representing
	 * the private member functions of the CLAS.
	 * 
	 * @param publicSharedFunctions
	 * A HashMap of String to LOLFunction representing
	 * the public SHARD functions of the CLAS.
	 * 
	 * @param privateSharedFunctions
	 * A HashMap of String to LOLFunction representing
	 * the private SHARD functions of the CLAS.
	 * 
	 * @param parentSource
	 * A String representing the name of the source
	 * file this CLAS was declared in.
	 */
	public LOLClass(String className, HashMap<String, ValueStruct> publicMemberVariables, HashMap<String, ValueStruct> privateMemberVariables, HashMap<String, ValueStruct> publicSharedVariables, HashMap<String, ValueStruct> privateSharedVariables, HashMap<String, LOLFunction> publicMemberFunctions, HashMap<String, LOLFunction> privateMemberFunctions, HashMap<String, LOLFunction> publicSharedFunctions, HashMap<String, LOLFunction> privateSharedFunctions, String parentSource) {
		this.className = className;
		this.publicMemberVariables = publicMemberVariables;
		this.privateMemberVariables = privateMemberVariables;
		this.publicSharedVariables = publicSharedVariables;
		this.privateSharedVariables = privateSharedVariables;
		this.publicMemberFunctions = publicMemberFunctions;
		this.privateMemberFunctions = privateMemberFunctions;
		this.publicSharedFunctions = publicSharedFunctions;
		this.privateSharedFunctions = privateSharedFunctions;
		this.parentSource = parentSource;
	}

	/**
	 * Gives the SHARD variable specified by the name, in
	 * the context of the specified LOLFunction.
	 * 
	 * If the specified variable is not found in the map
	 * of public SHARD variables and if the specified
	 * context is defined in the CLAS, then the variable
	 * is searched for in the private SHARD variables map.
	 * 
	 * @param name
	 * A String representing the name of the variable.
	 * 
	 * @param context
	 * A LOLFunction representing the location from which
	 * a request for the specified variable is initiated.
	 * 
	 * @return
	 * A ValueStruct corresponding to the specified variable,
	 * or null if not found in the given context's visibility
	 * permissions.
	 */
	public ValueStruct getSharedVariable(String name, LOLFunction context) {
		ValueStruct result = publicSharedVariables.get(name);

		if(result == null) {
			if(publicMemberFunctions.containsValue(context) || privateMemberFunctions.containsValue(context) || publicSharedFunctions.containsValue(context) || privateSharedFunctions.containsValue(context)) {
				result = privateSharedVariables.get(name);
			}
		}

		return result;
	}

	/**
	 * Gives the SHARD function specified by the name, in
	 * the context of the specified LOLFunction.
	 * 
	 * If the specified function is not found in the map
	 * of public SHARD functions and if the specified
	 * context is defined in the CLAS, then the function
	 * is searched for in the private SHARD function map.
	 * 
	 * @param name
	 * A String representing the name of the function.
	 * 
	 * @param context
	 * A LOLFunction representing the location from which
	 * a request for the specified function is initiated.
	 * 
	 * @return
	 * A LOLFunction corresponding to the specified function,
	 * or null if not found in the given context's visibility
	 * permissions.
	 */
	public LOLFunction getSharedFunction(String name, LOLFunction context) {
		LOLFunction result = publicSharedFunctions.get(name);

		if(result == null) {
			if(publicMemberFunctions.containsValue(context) || privateMemberFunctions.containsValue(context) || publicSharedFunctions.containsValue(context) || privateSharedFunctions.containsValue(context)) {
				result = privateSharedFunctions.get(name);
			}
		}

		return result;
	}
	
	/**
	 * Gives the member function specified by the name, in
	 * the context of the specified LOLFunction.
	 * 
	 * If the specified function is not found in the map
	 * of public member functions and if the specified
	 * context is defined in the CLAS, then the function
	 * is searched for in the private member function map.
	 * 
	 * @param name
	 * A String representing the name of the function.
	 * 
	 * @param context
	 * A LOLFunction representing the location from which
	 * a request for the specified function is initiated.
	 * 
	 * @return
	 * A LOLFunction corresponding to the specified function, 
	 * or null if not found in the given context's visibility
	 * permissions.
	 */
	public LOLFunction getMemberFunction(String name, LOLFunction context) {
		LOLFunction result = publicMemberFunctions.get(name);

		if(result == null) {
			if(publicMemberFunctions.containsValue(context) || privateMemberFunctions.containsValue(context)) {
				result = privateMemberFunctions.get(name);
			}
		}

		return result;
	}
	
	/**
	 * Checks if the specified LOLFunction is a member
	 * function of this CLAS.
	 * 
	 * @param context
	 * A LOLFunction to check for membership.
	 * 
	 * @return
	 * A boolean representing whether the specified LOLFunction
	 * is a member of this CLAS.
	 */
	public boolean isMemberFunction(LOLFunction context) {
		return publicMemberFunctions.containsValue(context) || privateMemberFunctions.containsValue(context);
	}
	
	/**
	 * Gives the name of this CLAS.
	 * 
	 * @return
	 * A String representing the CLAS name.
	 */
	public String getName() {
		return className;
	}

	/**
	 * Gives the parent source file of this CLAS.
	 * 
	 * @return
	 * A String representing the source name.
	 */
	public String getParentSource() {
		return parentSource;
	}

	/* (non-Javadoc)
	 * Checks for equality of objects. The other Object is
	 * equal to this instance of LOLClass if the other Object
	 * is an instance of LOLClass and if the CLAS and the
	 * parent source names are equivalent.
	 * 
	 * @see java.lang.Object#equals(java.lang.Object)
	 */
	@Override
	public boolean equals(Object o) {
		if(this == o) {
			return true;
		}

		if(!(o instanceof LOLClass)) {
			return false;
		}

		LOLClass rhs = (LOLClass)o;

		if(rhs.className.equals(className) && rhs.parentSource.equals(className)) {
			return true;
		}

		return false;
	}

	/* (non-Javadoc)
	 * Hashcode function. Adds the hashcodes of the CLAS
	 * name and source name, then bitshifts it to the right
	 * by 2.
	 * 
	 * @see java.lang.Object#hashCode()
	 */
	@Override
	public int hashCode() {
		return (className.hashCode() + parentSource.hashCode()) >> 2;
	}
	
	/**
	 * Gives a Collection of all the public member variables contained
	 * within the CLAS. Changing this Collection will cause the underlying
	 * map that contains the variables and their values to be altered
	 * as well.
	 * 
	 * @return
	 * A Collection of ValueStruct representing the public member
	 * variables defined by the CLAS.
	 * 
	 * @see java.util.HashMap#values()
	 */
	public Collection<ValueStruct> getPublicMemberVariables() {
		return publicMemberVariables.values();
	}
	
	/**
	 * Gives a Collection of all the private member variables contained
	 * within the CLAS. Changing this Collection will cause the underlying
	 * map that contains the variables and their values to be altered
	 * as well.
	 * 
	 * @return
	 * A Collection of ValueStruct representing the private member
	 * variables defined by the CLAS.
	 * 
	 * @see java.util.HashMap#values()
	 */
	public Collection<ValueStruct> getPrivateMemberVariables() {
		return privateMemberVariables.values();
	}
	
	/**
	 * Attempts to construct a LOLObject instance of this CLAS. If
	 * successful, returns the constructed object. If not, throws a
	 * LOLError.
	 * 
	 * @return
	 * A LOLObject representing the newly constructed instance of this
	 * CLAS.
	 * 
	 * @throws LOLError
	 * Throws a LOLError if instantiation fails. The LOLError contains a
	 * message of what went wrong.
	 * 
	 * @see org.objectivelol.lang.LOLObject
	 */
	public LOLObject constructInstance() throws LOLError {
		HashMap<String, ValueStruct> publicVars = new HashMap<String, ValueStruct>();
		HashMap<String, ValueStruct> privateVars = new HashMap<String, ValueStruct>();
		
		// copy the public member variables of the CLAS declaration to the LOLObject instance
		for(Iterator<Entry<String, ValueStruct>> i = publicMemberVariables.entrySet().iterator(); i.hasNext();) {
			Entry<String, ValueStruct> e = i.next();
			
			publicVars.put(e.getKey(), e.getValue().copy());
		}
		
		// copy the private member variables of the CLAS declaration to the LOLObject instance
		for(Iterator<Entry<String, ValueStruct>> i = privateMemberVariables.entrySet().iterator(); i.hasNext();) {
			Entry<String, ValueStruct> e = i.next();
			
			privateVars.put(e.getKey(), e.getValue().copy());
		}
		
		return new LOLObject(this, publicVars, privateVars);
	}
	
	/**
	 * Prepares the CLAS for use. Calls prepareFunction() in all of the SHARD and member
	 * functions. This is not a required function to call during runtime, since functions
	 * that have not yet been prepared will be automatically prepared before the first
	 * call. However, this function can be used to prepare functions for execution early,
	 * to reduce execution delays during runtime.
	 * 
	 * @throws LOLError
	 * Throws a LOLError if preparing the individual functions causes an error to arise.
	 * 
	 * @see org.objectivelol.lang.LOLFunction#prepareFunction()
	 */
	public void prepareClass() throws LOLError {
		for(LOLFunction func : publicSharedFunctions.values()) {
			func.prepareFunction();
		}
		
		for(LOLFunction func : privateSharedFunctions.values()) {
			func.prepareFunction();
		}
		
		for(LOLFunction func : publicMemberFunctions.values()) {
			func.prepareFunction();
		}
		
		for(LOLFunction func : privateMemberFunctions.values()) {
			func.prepareFunction();
		}
	}

}
