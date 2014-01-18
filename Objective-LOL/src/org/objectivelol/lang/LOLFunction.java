package org.objectivelol.lang;

import java.util.Collection;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.Map.Entry;

import org.objectivelol.lang.LOLValue.ValueStruct;
import org.objectivelol.vm.Expression;
import org.objectivelol.vm.Expression.Return;
import org.objectivelol.vm.Parser;

public class LOLFunction {

	private final String functionName;
	private final String returnType;
	private final LinkedHashMap<String, String> inputArguments;
	private final Boolean isShared;
	
	private final String parentClass;
	private final String parentSource;
	
	private String operations;
	private Expression expressions;
	
	public LOLFunction(String functionName, String returnType, LinkedHashMap<String, String> inputArguments, Boolean isShared, String parentClass, String parentSource, String operations) {
		this.functionName = functionName;
		this.returnType = returnType;
		this.inputArguments = inputArguments;
		this.isShared = isShared;
		this.parentClass = parentClass;
		this.parentSource = parentSource;
		this.operations = operations;
	}
	
	public String getArgumentType(String name) throws LOLError {
		String result = inputArguments.get(name);
		
		if(result == null) {
			throw new LOLError("Argument name not found");
		}
		
		return result;
	}
	
	public Collection<String> getArgumentTypes() {
		return inputArguments.values();
	}
	
	public String getParentClass() {
		return parentClass;
	}
	
	public String getParentSource() {
		return parentSource;
	}
	
	public String getName() {
		return functionName;
	}
	
	public void prepareFunction() throws LOLError {
		if(expressions == null) {
			expressions = Parser.parse(operations, this);
			operations = null;
		}
	}
	
	public final LOLValue execute(LOLObject owner, LOLValue ... args) throws LOLError {
		if(parentClass == null) {
			if(owner != null) {
				throw new LOLError("Invoking global functions does not require an owner object");
			}
		} else {
			if(!isShared) {
				if(owner == null) {
					throw new LOLError("Invoking a member function requires an owner object");
				}
			} else {
				if(owner != null) {
					throw new LOLError("Invoking a shared function does not require an owner object");
				}
			}
		}
		
		LinkedHashMap<String, ValueStruct> arguments = null;
		
		if((arguments = validateArguments(args)) == null) {
			throw new LOLError("Invalid number or types of arguments");
		}
		
		return run(owner, arguments);
	}
	
	protected LOLValue run(LOLObject owner, LinkedHashMap<String, ValueStruct> args) throws LOLError {
		try {
			if(expressions == null) {
				expressions = Parser.parse(operations, this);
				operations = null;
			}
			
			expressions.interpret(owner, this, args);
			return LOLNothing.NOTHIN;
		} catch(Return e) {
			return e.getValue();
		}
	}
	
	protected LinkedHashMap<String, ValueStruct> validateArguments(LOLValue ... args) {
		if(args.length != inputArguments.size()) {
			return null;
		}
		
		LinkedHashMap<String, ValueStruct> result = new LinkedHashMap<String, ValueStruct>();
		
		int counter = 0;
		for(Iterator<Entry<String, String>> i = inputArguments.entrySet().iterator(); i.hasNext();) {
			Entry<String, String> e = i.next();
			
			if(!e.getValue().equals(args[counter].getTypeName())) {
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
	
	public String getReturnType() {
		return returnType;
	}
	
	@Override
	public boolean equals(Object o) {
		if(this == o) {
			return true;
		}
		
		if(!(o instanceof LOLFunction)) {
			return false;
		}
		
		LOLFunction rhs = (LOLFunction)o;
		
		if(rhs.functionName.equals(functionName) && rhs.parentClass.equals(parentClass) && rhs.inputArguments.equals(inputArguments)) {
			return true;
		}
		
		return false;
	}
	
	@Override
	public int hashCode() {
		return (functionName.hashCode() + parentClass.hashCode() + inputArguments.hashCode()) >> 2;
	}

	public Boolean isShared() {
		return isShared;
	}
	
	public int numArgs() {
		return inputArguments.size();
	}

}
