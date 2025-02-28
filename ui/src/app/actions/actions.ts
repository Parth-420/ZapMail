"use server";

import { parseEmail } from "@/hooks/emailParser";
import { db } from "@/lib/db";

export async function searchEmails(rcptQuery: string) {
  try {
    console.log("Searching for emails with query:", rcptQuery);
    
    // Format the search query to match the database format with angle brackets
    const formattedQuery = `<${rcptQuery}@zapmail.parth.lol>`;
    console.log("Formatted query:", formattedQuery);
    
    const result = await db.query(
      `SELECT id, username as mail_from, recipient as rcpt_to, raw_data as data, received_at as date 
       FROM emails 
       WHERE recipient = $1
       ORDER BY received_at DESC`,
      [formattedQuery]
    );
    
    console.log("Raw query results:", result.rows);
    
    const output = [];
    for (const i of result.rows) {
      console.log("Processing row:", i);
      try {
        const parsedEmail = await parseEmail(i.data);
        console.log("Parsed email:", parsedEmail);
        output.push({
          id: i.id,
          date: i.date,
          mail_from: i.mail_from,
          rcpt_to: i.rcpt_to,
          data: parsedEmail,
        });
      } catch (parseError) {
        console.error("Error parsing email:", parseError);
        // Continue with next email if one fails to parse
      }
    }
    
    console.log("Final output:", output);
    return output;
  } catch (error) {
    console.error("Error fetching emails", error);
    throw new Error("Error fetching emails");
  }
}
