package com.application.cs260final;

import androidx.appcompat.app.AppCompatActivity;

import android.app.AlertDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.Toast;

public class MakeNewPost extends AppCompatActivity {

    Button postButton;
    ImageView xButton;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_home_screen_make_a_post);

        postButton = findViewById(R.id.post);
        postButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                Toast.makeText(MakeNewPost.this, "You made a new post!", Toast.LENGTH_LONG).show();
                //what else should we do to indicate that the user made a new post?
                //this just takes the user back to the homepage after clicking post lol
                Intent intent = new Intent(MakeNewPost.this, HomeScreenActivity.class);
                startActivity(intent);
            }
        });

        xButton = findViewById(R.id.imageXout);
        xButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                //maybe have an if statement: (if text input !null, have some option to save draft? idk)
                alertDialog();
            }
        });
    }

    private void alertDialog() {
        AlertDialog.Builder dialog=new AlertDialog.Builder(this);
        dialog.setTitle("Do you want to discard your post?");
        dialog.setPositiveButton("YES",
                new DialogInterface.OnClickListener() {
                    public void onClick(DialogInterface dialog,
                                        int which) {
                        Toast.makeText(getApplicationContext(),"Your post has been discarded",Toast.LENGTH_LONG).show();
                        Intent intent = new Intent(MakeNewPost.this, HomeScreenActivity.class);
                        startActivity(intent);
                    }
                });
        dialog.setNegativeButton("NO",new DialogInterface.OnClickListener() {
            @Override
            public void onClick(DialogInterface dialog, int which) {
                Toast.makeText(getApplicationContext(),"Continue writing!",Toast.LENGTH_LONG).show();
            }
        });
        AlertDialog alertDialog=dialog.create();
        alertDialog.show();
    }


}
