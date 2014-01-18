package org.objectivelol.libs;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;

import org.objectivelol.lang.LOLError;
import org.objectivelol.lang.LOLNative;
import org.objectivelol.lang.LOLNothing;
import org.objectivelol.lang.LOLString;
import org.objectivelol.lang.LOLValue;

public class STDIO extends LOLNative {
	
	private static BufferedReader br = new BufferedReader(new InputStreamReader(System.in));

	public static LOLNothing COMPLAIN(LOLString arg) {
		System.err.println(arg.toString());
		return LOLNothing.NOTHIN;
	}
	
	public static LOLString GIMMEH() {
		try {
			return (LOLString)LOLValue.valueOf(br.readLine()).cast(LOLString.TYPE_NAME);
		} catch(IOException e) {
			return null;
		} catch(LOLError e) {
			return null;
		}
	}
	
	public static LOLNothing VISIBLE(LOLString arg) {
		System.out.println(arg.toString());
		return LOLNothing.NOTHIN;
	}
	
}
