package org.objectivelol.lang;

import java.util.Collection;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map.Entry;

import org.objectivelol.vm.ValueStruct;

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

	public ValueStruct getSharedVariable(String name, LOLFunction context) {
		ValueStruct result = publicSharedVariables.get(name);

		if(result == null) {
			if(publicMemberFunctions.containsValue(context) || privateMemberFunctions.containsValue(context) || publicSharedFunctions.containsValue(context) || privateSharedFunctions.containsValue(context)) {
				result = privateSharedVariables.get(name);
			}
		}

		return result;
	}

	public LOLFunction getSharedFunction(String name, LOLFunction context) {
		LOLFunction result = publicSharedFunctions.get(name);

		if(result == null) {
			if(publicMemberFunctions.containsValue(context) || privateMemberFunctions.containsValue(context) || publicSharedFunctions.containsValue(context) || privateSharedFunctions.containsValue(context)) {
				result = privateSharedFunctions.get(name);
			}
		}

		return result;
	}
	
	public LOLFunction getMemberFunction(String name, LOLFunction context) {
		LOLFunction result = publicMemberFunctions.get(name);

		if(result == null) {
			if(publicMemberFunctions.containsValue(context) || privateMemberFunctions.containsValue(context)) {
				result = privateMemberFunctions.get(name);
			}
		}

		return result;
	}
	
	public boolean isMemberFunction(LOLFunction context) {
		return publicMemberFunctions.containsValue(context) || privateMemberFunctions.containsValue(context);
	}
	
	public String getName() {
		return className;
	}

	public String getParentSource() {
		return parentSource;
	}

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

	@Override
	public int hashCode() {
		return (className.hashCode() + parentSource.hashCode()) >> 2;
	}
	
	public Collection<ValueStruct> getPublicMemberVariables() {
		return publicMemberVariables.values();
	}
	
	public Collection<ValueStruct> getPrivateMemberVariables() {
		return privateMemberVariables.values();
	}
	
	public LOLObject constructInstance() throws LOLError {
		HashMap<String, ValueStruct> publicVars = new HashMap<String, ValueStruct>();
		HashMap<String, ValueStruct> privateVars = new HashMap<String, ValueStruct>();
		
		for(Iterator<Entry<String, ValueStruct>> i = publicMemberVariables.entrySet().iterator(); i.hasNext();) {
			Entry<String, ValueStruct> e = i.next();
			
			publicVars.put(e.getKey(), e.getValue().copy());
		}
		
		for(Iterator<Entry<String, ValueStruct>> i = privateMemberVariables.entrySet().iterator(); i.hasNext();) {
			Entry<String, ValueStruct> e = i.next();
			
			privateVars.put(e.getKey(), e.getValue().copy());
		}
		
		return new LOLObject(this, publicVars, privateVars);
	}
	
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
