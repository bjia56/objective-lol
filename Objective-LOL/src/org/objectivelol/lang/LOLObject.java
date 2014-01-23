package org.objectivelol.lang;

import java.util.HashMap;
import java.util.Map.Entry;

import org.objectivelol.vm.ValueStruct;

public class LOLObject extends LOLValue {

	private final LOLClass objectType;

	private final HashMap<String, ValueStruct> publicMemberVariables;
	private final HashMap<String, ValueStruct> privateMemberVariables;

	public LOLObject(LOLClass objectType, HashMap<String, ValueStruct> publicMemberVariables, HashMap<String, ValueStruct> privateMemberVariables) {
		this.objectType = objectType;
		this.publicMemberVariables = publicMemberVariables;
		this.privateMemberVariables = privateMemberVariables;
	}

	@Override
	public LOLValue cast(String type) throws LOLError {
		if(type.equals(objectType.getName())) {
			return this;
		}

		throw new LOLError("Cannot cast to the specified type");
	}

	@Override
	public String getTypeName() {
		return objectType.getName();
	}

	public LOLFunction getFunction(String name, LOLFunction context) throws LOLError {
		LOLFunction result = objectType.getMemberFunction(name, context);

		if(result == null) {
			result = objectType.getSharedFunction(name, context);
		}
		
		return result;
	}

	public ValueStruct getVariable(String name, LOLFunction context) throws LOLError {
		ValueStruct result = objectType.getSharedVariable(name, context);
		
		if(result == null) {
			result = publicMemberVariables.get(name);
			
			if(result == null) {
				if(objectType.isMemberFunction(context)) {
					result = privateMemberVariables.get(name);
				}
			}
		}
		
		return result;
	}

	@Override
	public LOLBoolean equalTo(LOLValue other) throws LOLError {
		LOLObject lo = (LOLObject)other.cast(objectType.getName());
		
		return (lo.privateMemberVariables.equals(privateMemberVariables) && lo.publicMemberVariables.equals(publicMemberVariables) ? LOLBoolean.YEZ : LOLBoolean.NO);
	}

	@Override
	public LOLValue copy() throws LOLError {
		HashMap<String, ValueStruct> publicMemberVariables = new HashMap<String, ValueStruct>();
		HashMap<String, ValueStruct> privateMemberVariables = new HashMap<String, ValueStruct>();
		
		for(Entry<String, ValueStruct> e : this.publicMemberVariables.entrySet()) {
			publicMemberVariables.put(e.getKey(), e.getValue().copy());
		}
		
		for(Entry<String, ValueStruct> e : this.privateMemberVariables.entrySet()) {
			privateMemberVariables.put(e.getKey(), e.getValue().copy());
		}
		
		return new LOLObject(objectType, publicMemberVariables, privateMemberVariables);
	}

}
