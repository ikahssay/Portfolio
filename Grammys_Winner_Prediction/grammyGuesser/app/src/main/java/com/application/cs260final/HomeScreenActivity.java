package com.application.cs260final;

import android.content.Intent;
import android.graphics.drawable.Drawable;
import android.os.Bundle;
import android.view.MenuItem;
import android.view.View;
import android.widget.ImageButton;
import android.widget.ImageView;
import android.widget.TextView;
import android.widget.Toast;


import androidx.annotation.NonNull;
import androidx.appcompat.app.AppCompatActivity;

import com.google.android.material.bottomnavigation.BottomNavigationView;


public class HomeScreenActivity extends AppCompatActivity{
    ImageView make_post_button;
    ImageView click_on_post_button;
    ImageView upvote1;
    ImageView downvote1;
    ImageView upvote2;
    ImageView downvote2;
    ImageView search_button;
    ImageView predictions_button;
    ImageView notification_button;
    ImageView profile_button;
    TextView textupvote;
    TextView textdownvote;

    TextView textupvote2;
    TextView textdownvote2;

    @Override
    protected void onCreate(Bundle savedInstanceState){
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_home_screen);

        //1. CREATE AN ACTION LINKED TO THE "MAKE A POST" BUTTON
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        make_post_button = (ImageView) findViewById(R.id.floatingBtnFloatingactionbutton);
        make_post_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(HomeScreenActivity.this, "Make a New Post", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(HomeScreenActivity.this, MakeNewPost.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //2. CREATE AN ACTION LINKED TO THE CLICK ON POST BUTTON
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        click_on_post_button = (ImageView) findViewById(R.id.imageForwardButton);
        click_on_post_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(HomeScreenActivity.this, "Details of Post", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(HomeScreenActivity.this, ShowDetailedPostScreen.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //3. CREATE AN ACTION LINKED TO THE DOWNVOTING BUTTON -- make this work for both downvoting and upvoting
        // somehow update the number of up/down votes on the post?
        // also it doesn't save once you leave the page, only saves when you use the back button

        //on Char's post
        upvote1 = (ImageView) findViewById(R.id.imageUpvote20);
        textupvote = (TextView)  findViewById(R.id.txt20);


        downvote1 = (ImageView) findViewById(R.id.imageDownvote4);
        textdownvote =(TextView) findViewById(R.id.txt4) ;

        upvote1.setOnClickListener(new View.OnClickListener(){
            public void onClick(View v) {
                upvote1.setActivated(!upvote1.isActivated());
                textupvote.setText("21");

            }
        });
        downvote1.setOnClickListener(new View.OnClickListener(){
            public void onClick(View v) {
                downvote1.setActivated(!downvote1.isActivated());
                textdownvote.setText("5");
            }
        });

        //on John's post
        upvote2 = (ImageView) findViewById(R.id.imageUpvote21);
        textupvote2 = (TextView) findViewById(R.id.txt2);
        upvote2.setOnClickListener(new View.OnClickListener(){
            public void onClick(View v) {
                upvote2.setActivated(!upvote2.isActivated());
                textupvote2.setText("3");

            }
        });
        downvote2 = (ImageView) findViewById(R.id.imageDownvote5);
        textdownvote2 = (TextView)findViewById(R.id.txt100);
        downvote2.setOnClickListener(new View.OnClickListener(){
            public void onClick(View v) {
                downvote2.setActivated(!downvote2.isActivated());
                textdownvote2.setText("101");
            }
        });

        // ************************************* MENU BUTTON FUNCTIONALITIES BELOW **************************************** //

        //4. CREATE AN ACTION LINKED TO THE SEARCH BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        search_button = (ImageView) findViewById(R.id.imageSearch);
        search_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(HomeScreenActivity.this, "Going to Search Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(HomeScreenActivity.this, Search.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //5. CREATE AN ACTION LINKED TO THE PREDITCION BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        predictions_button = (ImageView) findViewById(R.id.imageVote);
        predictions_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(HomeScreenActivity.this, "Going to Predictions Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(HomeScreenActivity.this, Predictions.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //6. CREATE AN ACTION LINKED TO THE NOTIFICATIONS BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        notification_button = (ImageView) findViewById(R.id.imageBell);
        notification_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(HomeScreenActivity.this, "Going to Notifications Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(HomeScreenActivity.this, Notifications.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //7. CREATE AN ACTION LINKED TO THE PROFILE BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        profile_button = (ImageView) findViewById(R.id.imageProfileButton);
        profile_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(HomeScreenActivity.this, "Going to Profile Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(HomeScreenActivity.this, Profile.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

    }

}
