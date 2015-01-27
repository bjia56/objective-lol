package org.objectivelol.lang;

import java.util.Collection;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.Map.Entry;

import org.objectivelol.vm.Expression;
import org.objectivelol.vm.Expression.Return;
import org.objectivelol.vm.Parser;
import org.objectivelol.vm.ValueStruct;

/**
 * Class to represent a function in Objective-LOL.
 * 
 * @author Brett Jia
 */
public class LOLFunction {

	private final String functionName;
	private final String returnType;
	private final LinkedHashMap<String, String> inputArguments;
	private final Boolean isShared;
	
	private final String parentClass;
	private final String parentSource;
	
	private String operations;
	private Expression expressions;
	
	/**
	 * Constructor for the LOLFunction class.
	 * 
	 * @param functionName
	 * A String representing the name of the function.
	 * 
	 * @param returnType
	 * A String representing the return type of the function.
	 * This value can be null if the function does not
	 * return anything; in such a case, the LOLFunction
	 * will return LOLNothing at runtime.
	 * 
	 * @param inputArguments
	 * A LinkedHashMap of String to String representing a mapping
	 * of the argument name to the argument type. This must be
	 * a map representation that guarantees the order of the
	 * elements, since the order of elements in this map must
	 * match the order of arguments upon declaration in a
	 * source file.
	 * 
	 * @param isShared
	 * A Boolean representing whether this function is SHARD or not.
	 * If a null value is passed in, then the function is global.
	 * 
	 * @param parentClass
	 * A String representing the name of the parent class of this
	 * function, if applicable.
	 * 
	 * @param parentSource
	 * A String representing the name of the source file that contains
	 * this function.
	 * 
	 * @param operations
	 * A String representing all of the operations that this function
	 * contains. This is essentially the raw text of the function's
	 * body.
	 * 
	 * @see java.util.LinkedHashMap
	 */
	public LOLFunction(String functionName, String returnType, LinkedHashMap<String, String> inputArguments, Boolean isShared, String parentClass, String parentSource, String operations) {
		this.functionName = functionName;
		if(returnType == null) {
			this.returnType = LOLNothing.TYPE_NAME;
		} else {
			this.returnType = returnType;
		}
		this.inputArguments = inputArguments;
		this.isShared = isShared;
		this.parentClass = parentClass;
		this.parentSource = parentSource;
		this.operations = operations;
	}
	
	/**
	 * Gives the type of a parameter argument of this LOLFunction
	 * instance.
	 * 
	 * @param name
	 * A String representing the name of the argument to retrieve.
	 * 
	 * @return
	 * A String representing the given argument's type.
	 * 
	 * @throws LOLError
	 * Throws a LOLError if the argument is not found.
	 */
	public String getArgumentType(String name) throws LOLError {
		String result = inputArguments.get(name);
		
		if(result == null) {
			throw new LOLError("Argument name not found");
		}
		
		return result;
	}
	
	/**
	 * Gives a Collection of all the argument types defined for this
	 * function. Changing this Collection will cause the underlying
	 * map that contains the arguments and their type representations
	 * to be altered as well.
	 * 
	 * @return
	 * A Collection of String representing the argument types  of the
	 * function as defined by the source code.
	 * 
	 * @see java.util.HashMap#values()
	 */
	public Collection<String> getArgumentTypes() {
		return inputArguments.values();
	}
	
	/**
	 * Gives a String representing the name of the parent class of this
	 * function. If the function is global, a null value can be returned.
	 * 
	 * @return
	 * A String representing the name of the parent class of this function,
	 * null if the function is global.
	 */
	public String getParentClass() {
		return parentClass;
	}
	
	/**
	 * Gives a String representing the name of the source that this function
	 * can be found in.
	 * 
	 * @return
	 * A String representing the name of the source this function is contained
	 * in.
	 */
	public String getParentSource() {
		return parentSource;
	}
	
	/**
	 * Gives a String representing the name of this function.
	 * 
	 * @return
	 * A String representing the name of this function.
	 */
	public String getName() {
		return functionName;
	}
	
	/**
	 * Prepares this function by parsing all the operations contained in the function
	 * declaration if the operations have not already been parsed. The String representation
	 * of these operations will be cleared to save memory.
	 * 
	 * @throws LOLError
	 * Throws a LOLError if an error occurs in parsing.
	 * 
	 * @see org.objectivelol.vm.Parser
	 */
	public void prepareFunction() throws LOLError {
		if(expressions == null) {
			expressions = Parser.parse(operations, this);
			operations = null;
		}
	}
	
	/**
	 * Runs this function by executing all the operations contained in the function
	 * sequentially. Invoking a global function requires that an owner object is
	 * not specified. Invoking a member function of a CLAS requires an owner object,
	 * while invoking a SHARD function of a CLAS requires that an owner object is not
	 * specified.
	 * 
	 * @param owner
	 * A LOLObject specifying which object owns this function. If this function is
	 * global or SHARD, then this argument must be null.
	 * 
	 * @param args
	 * A LOLValue varargs specifying the arguments that are passed into this function.
	 * For the function to execute, the number of arguments passed must match the number
	 * of arguments present upon declaration of the function, and each argument must
	 * be able to be cast into the declared argument's type. Not including this value
	 * or passing in null or an empty array of LOLValue will act as calling a function
	 * without arguments.
	 * 
	 * @return
	 * A LOLValue representing the result of the execution. A LOLNothing is returned if
	 * the function does not have a return type.
	 * 
	 * @throws LOLError
	 * Throws a LOLError if the function is unable to be executed due to argument
	 * mismatches, or if execution produces an error. An execution error can be caused
	 * by errors in parsing, errors in runtime execution, or mismatched return types.
	 */
	public final LOLValue execute(LOLObject owner, LOLValue ... args) throws LOLError {
		// check if function is global
		if(parentClass == null) {
			// check if an owner object is specified
			if(owner != null) {
				throw new LOLError("Invoking global functions requires absence of owner object");
			}
		} else {
			if(!isShared) {
				// check if an owner object is specified for a member function invocation
				if(owner == null) {
					throw new LOLError("Invoking a member function requires an owner object");
				}
				// check if the owner's type is the parent class of this function
				if(!owner.getTypeName().equals(parentClass)) {
					throw new LOLError("Member function invocation requires the function to be a member of the specified object");
				}
			} else {
				// check if an owner object is specified for a SHARD function invocation
				if(owner != null) {
					throw new LOLError("Invoking a shared function does not require an owner object");
				}
			}
		}
		
		// prevents NullPointerException from being raised when executing function's expressions
		if(args == null) {
			args = new LOLValue[0];
		}
		
		LinkedHashMap<String, ValueStruct> arguments = null;
		
		// check for valid input arguments
		if((arguments = validateArguments(args)) == null) {
			throw new LOLError("Invalid number or types of arguments");
		}
		
		return run(owner, arguments);
	}
	
	/**
	 * Executes the actual expressions in the function body, assuming all error checking
	 * of argument and type compatibility has already been completed.
	 * 
	 * @param owner
	 * A LOLObject representing the owner of this function, if it is a member function.
	 * 
	 * @param args
	 * A LinkedHashMap of String to ValueStruct representing the input arguments of the
	 * function. This is assumed to be accurate in number and types, following the function
	 * declaration in the source file.
	 * 
	 * @return
	 * A LOLValue representing the result of execution. Will return a LOLNothing if the
	 * function does not have a return type.
	 * 
	 * @throws LOLError
	 * Throws a LOLError if parsing the statements raises an error (if the statements have
	 * not yet been parsed), if executing the commands raises an error, or if the function
	 * return value's type does not match the specified function return type. If a value is
	 * returned but no value was expected, then a LOLError is thrown. If a value is not
	 * returned but a value was expected, then a LOLError is also thrown.
	 */
	protected LOLValue run(LOLObject owner, LinkedHashMap<String, ValueStruct> args) throws LOLError {
		try {
			if(expressions == null) {
				expressions = Parser.parse(operations, this);
				operations = null;
			}
			
			expressions.interpret(owner, this, args);
			
			if(!LOLNothing.TYPE_NAME.equals(returnType)) {
				throw new LOLError("Execution of function with return type ended without encountering return statement");
			}
			
			// functions with no return type are to return a LOLNothing
			return LOLNothing.NOTHIN;
		} catch(Return e) { // returned values are thrown to easily cease execution and return the value needed
			LOLValue v = e.getValue();
			
			// check if the return actually produced a value, or was an early exit from a function with no return type
			if(v == null) {
				if(!LOLNothing.TYPE_NAME.equals(returnType)) {
					throw new LOLError("Execution of function with return type ended without encountering return statement");
				}
				
				return LOLNothing.NOTHIN;
			} else {
				if(LOLNothing.TYPE_NAME.equals(returnType)) {
					throw new LOLError("Execution of function without return type ended with return statement");
				}
				
				return e.getValue();
			}
		}
	}
	
	/**
	 * Validates the arguments passed into the function and attempts to cast them to
	 * the specified types.
	 * 
	 * @param args
	 * A LOLValue varargs representing the arguments passed in.
	 * 
	 * @return
	 * A LinkedHashMap of String to ValueStruct representing the initialized and copied
	 * argument values of the function execution. If an issue arises during argument
	 * validation, such as argument number or type mismatch, then returns null.
	 */
	protected LinkedHashMap<String, ValueStruct> validateArguments(LOLValue ... args) {
		if(args.length != inputArguments.size()) {
			return null;
		}
		
		LinkedHashMap<String, ValueStruct> result = new LinkedHashMap<String, ValueStruct>();
		
		int counter = 0;
		for(Iterator<Entry<String, String>> i = inputArguments.entrySet().iterator(); i.hasNext();) {
			Entry<String, String> e = i.next();
			
			if(!e.getValue().equals(args[counter].getTypeName())) {
				// if the type of the argument passed in is not the same as the specified type, try to cast it to the accepted type
				try {
					args[counter] = args[counter].cast(e.getValue());
				} catch(LOLError l) {
					return null;
				}
			}
			
			result.put(e.getKey(), new ValueStruct(e.getValue(), args[counter++], false));
		}
		
		return result;
	}
	
	/**
	 * Gives a String representing the return type of this function.
	 * 
	 * @return
	 * A String representing the return type of this function.
	 */
	public String getReturnType() {
		return returnType;
	}
	
	/* (non-Javadoc)
	 * Checks for equality of objects. The other Object is
	 * equal to this instance of LOLFunction if the other Object
	 * is an instance of LOLFunction and contains the same
	 * function name and parent class name.
	 * 
	 * @see java.lang.Object#equals(java.lang.Object)
	 */
	@Override
	public boolean equals(Object o) {
		if(this == o) {
			return true;
		}
		
		if(!(o instanceof LOLFunction)) {
			return false;
		}
		
		LOLFunction rhs = (LOLFunction)o;
		
		if(rhs.functionName.equals(functionName) && rhs.parentClass.equals(parentClass)) {
			return true;
		}
		
		return false;
	}
	
	/* (non-Javadoc)
	 * Hashcode function. Adds the hashcodes of the function name and
	 * parent class name, then bitshifts it to the right by 2.
	 * 
	 * @see java.lang.Object#hashCode()
	 */
	@Override
	public int hashCode() {
		return (functionName.hashCode() + parentClass.hashCode()) >> 2;
	}

	/**
	 * Gives a Boolean denoting whether the function is SHARD or not.
	 * If the function is global, a null value is returned.
	 * 
	 * @return
	 * A Boolean representing whether the function is SHARD, or null
	 * if it is global.
	 */
	public Boolean isShared() {
		return isShared;
	}
	
	/**
	 * Gives an int representing the size of the argument list of this
	 * function.
	 * 
	 * @return
	 * An int representing the number of arguments that this function
	 * takes in.
	 */
	public int numArgs() {
		return inputArguments.size();
	}

}
