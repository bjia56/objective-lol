package org.objectivelol.libs;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import org.objectivelol.lang.LOLNative;
import org.objectivelol.lang.LOLString;
import org.objectivelol.lang.LOLValue;

public class STDIO extends LOLNative {
	
	private static BufferedReader br = new BufferedReader(new InputStreamReader(System.in));

	public static void COMPLAIN(LOLString arg) {
		System.err.println(arg.toString());
	}
	
	public static LOLString GIMMEH() {
		try {
			return (LOLString)LOLValue.valueOf(br.readLine());
		} catch (IOException e) {
			return null;
		}
	}
	
	public static void VISIBLE(LOLString arg) {
		System.out.println(arg.toString());
	}
	
}
