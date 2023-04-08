# valueIterationAgents.py
# -----------------------
# Licensing Information:  You are free to use or extend these projects for
# educational purposes provided that (1) you do not distribute or publish
# solutions, (2) you retain this notice, and (3) you provide clear
# attribution to UC Berkeley, including a link to http://ai.berkeley.edu.
#
# Attribution Information: The Pacman AI projects were developed at UC Berkeley.
# The core projects and autograders were primarily created by John DeNero
# (denero@cs.berkeley.edu) and Dan Klein (klein@cs.berkeley.edu).
# Student side autograding was added by Brad Miller, Nick Hay, and
# Pieter Abbeel (pabbeel@cs.berkeley.edu).


# valueIterationAgents.py
# -----------------------
# Licensing Information:  You are free to use or extend these projects for
# educational purposes provided that (1) you do not distribute or publish
# solutions, (2) you retain this notice, and (3) you provide clear
# attribution to UC Berkeley, including a link to http://ai.berkeley.edu.
#
# Attribution Information: The Pacman AI projects were developed at UC Berkeley.
# The core projects and autograders were primarily created by John DeNero
# (denero@cs.berkeley.edu) and Dan Klein (klein@cs.berkeley.edu).
# Student side autograding was added by Brad Miller, Nick Hay, and
# Pieter Abbeel (pabbeel@cs.berkeley.edu).


import mdp, util

from learningAgents import ValueEstimationAgent
import collections

class ValueIterationAgent(ValueEstimationAgent):
    """
        * Please read learningAgents.py before reading this.*

        A ValueIterationAgent takes a Markov decision process
        (see mdp.py) on initialization and runs value iteration
        for a given number of iterations using the supplied
        discount factor.
    """
    def __init__(self, mdp, discount = 0.9, iterations = 100):
        """
          Your value iteration agent should take an mdp on
          construction, run the indicated number of iterations
          and then act according to the resulting policy.

          Some useful mdp methods you will use:
              mdp.getStates()
              mdp.getPossibleActions(state)
              mdp.getTransitionStatesAndProbs(state, action)
              mdp.getReward(state, action, nextState)
              mdp.isTerminal(state)
        """
        self.mdp = mdp
        self.discount = discount
        self.iterations = iterations
        self.values = util.Counter() # A Counter is a dict with default 0
        self.runValueIteration()

    def runValueIteration(self):
        # Write value iteration code here
        "*** YOUR CODE HERE ***"

        #Gives a list of all the states
        allStates = self.mdp.getStates()
        k = self.iterations
        discount = self.discount

        #Initially gives an empty dictionary
        #Initial iteration -> all states are initialized to 0
        V_k = util.Counter()
        V_k_1 = util.Counter()
        policy = util.Counter()

        #print("State values: " + str(V_k_1))

        #Already initialized all values when counter = 0
        counter = 1

        while counter <= k:
            #Look at each state and compute the values for each iteration k
            for i in range(0, len(allStates)):

                #Get all possible actions EACH STATE can take
                    #Ex: Get all legal actions from state 1
                legal_actions = self.mdp.getPossibleActions(allStates[i])
                Q_value_for_each_action = []

                #Look at each action each state can take, and compute the Q values
                    #Ex: Get all transition states from each legal action state 1 can take
                for eachAction in legal_actions:

                    # In transition_states, it has a TUPLE of successor states and the probability to get to it.
                        #Ex: Get all transition states from each legal action state 1 can take
                    transition_states = self.mdp.getTransitionStatesAndProbs(allStates[i], eachAction)

                    sum = 0

                    #Calculates the Q value of each state,action pair
                    for eachSuccessor, probability in transition_states:
                        reward = self.mdp.getReward(allStates[i], eachAction, eachSuccessor)
                        sum += probability * (reward  + (discount * V_k_1[eachSuccessor]))

                    Q_value_for_each_action.append(sum)

                if len(Q_value_for_each_action) == 0:
                    V_k[allStates[i]] = 0
                    policy[allStates[i]] = None
                else:
                    #Value iteration takes the best Q-value (which corresponds to the best action to take)
                    V_k[allStates[i]] = max(Q_value_for_each_action)

                    #Gets the highest value's index in the Q_value list
                    #best_index = np.argmax(Q_value_for_each_action)
                    best_index = Q_value_for_each_action.index(max(Q_value_for_each_action))

                    #Uses best_index to get the best action
                    best_action = legal_actions[best_index]

                    #Policy is a DICTIONARY-> For each state, there will be a corresponding best action to take.
                    policy[allStates[i]] = best_action

            counter = counter + 1

            #Checks if values converge:
            if (V_k_1 == V_k):
                break

            #Updates prev iteration (V_k-1) before looping again
            V_k_1 = V_k.copy()

        self.values = V_k.copy()

        return policy


    def getValue(self, state):
        """
          Return the value of the state (computed in __init__).
        """
        return self.values[state]


    def computeQValueFromValues(self, state, action):
        """
          Compute the Q-value of action in state from the
          value function stored in self.values.
        """
        "*** YOUR CODE HERE ***"

        V_k_1 = self.values
        discount = self.discount

        # In transition_states, it has a tuple of successor states and the probability to get to it.
            #Ex: Get all transition states from each legal action state 1 can take
        transition_states = self.mdp.getTransitionStatesAndProbs(state, action)

        Q_value = 0
        for eachSuccessor, probability in transition_states:
            reward = self.mdp.getReward(state, action, eachSuccessor)
            Q_value += probability * (reward  + (discount * V_k_1[eachSuccessor]))

        return Q_value

    def computeActionFromValues(self, state):
        """
          The policy is the best action from the given state
          according to the values currently stored in self.values.

          You may break ties any way you see fit.  Note that if
          there are no legal actions, which is the case at the
          terminal state, you should return None.
        """
        "*** YOUR CODE HERE ***"

        V_k_1 = self.values
        discount = self.discount
        legal_actions = self.mdp.getPossibleActions(state)
        Q_value_for_each_action = []

        #if self.mdp.isTerminal(state):
        if len(legal_actions) == 0:
            return None
        else:
            #Look at each action each state can take, and compute the Q values
                #Ex: Get all transition states from each legal action state 1 can take
            for eachAction in legal_actions:
                # In transition_states, it has a tuple of successor states and the probability to get to it.
                    #Ex: Get all transition states from each legal action state 1 can take
                transition_states = self.mdp.getTransitionStatesAndProbs(state, eachAction)

                sum = 0.0
                for eachSuccessor, probability in transition_states:
                    reward = self.mdp.getReward(state, eachAction, eachSuccessor)
                    sum += probability * (reward  + (discount * float(V_k_1[eachSuccessor])))

                Q_value_for_each_action.append(sum)

            #Gets the highest value's index in the Q_value list
            #best_index = np.argmax(Q_value_for_each_action)

            best_index = Q_value_for_each_action.index(max(Q_value_for_each_action))
            #Uses best_index to get the best action
            best_action = legal_actions[best_index]

            return best_action


    def getPolicy(self, state):
        return self.computeActionFromValues(state)

    def getAction(self, state):
        "Returns the policy at the state (no exploration)."
        return self.computeActionFromValues(state)

    def getQValue(self, state, action):
        return self.computeQValueFromValues(state, action)

class AsynchronousValueIterationAgent(ValueIterationAgent):
    """
        * Please read learningAgents.py before reading this.*

        An AsynchronousValueIterationAgent takes a Markov decision process
        (see mdp.py) on initialization and runs cyclic value iteration
        for a given number of iterations using the supplied
        discount factor.
    """
    def __init__(self, mdp, discount = 0.9, iterations = 1000):
        """
          Your cyclic value iteration agent should take an mdp on
          construction, run the indicated number of iterations,
          and then act according to the resulting policy. Each iteration
          updates the value of only one state, which cycles through
          the states list. If the chosen state is terminal, nothing
          happens in that iteration.

          Some useful mdp methods you will use:
              mdp.getStates()
              mdp.getPossibleActions(state)
              mdp.getTransitionStatesAndProbs(state, action)
              mdp.getReward(state)
              mdp.isTerminal(state)
        """
        ValueIterationAgent.__init__(self, mdp, discount, iterations)

    def runValueIteration(self):
        "*** YOUR CODE HERE ***"
        #Gives a list of all the states
        allStates = self.mdp.getStates()
        k = self.iterations
        discount = self.discount
    
        #Initially gives an empty dictionary
        #Initial iteration -> all states are initialized to 0
        V_k = util.Counter()
        #V_k_1 = util.Counter()
        policy = util.Counter()
         
        #print("State values: " + str(V_k_1))

        #Already initialized all values when counter = 0
        counter = 1
        i = 0
        
        while counter <= k:
           
            #Get all possible actions EACH STATE can take
                #Ex: Get all legal actions from state 1
            legal_actions = self.mdp.getPossibleActions(allStates[i])
            Q_value_for_each_action = []
            
            #Look at each action each state can take, and compute the Q values
                #Ex: Get all transition states from each legal action state 1 can take
            for eachAction in legal_actions:
            
                # In transition_states, it has a TUPLE of successor states and the probability to get to it.
                    #Ex: Get all transition states from each legal action state 1 can take
                transition_states = self.mdp.getTransitionStatesAndProbs(allStates[i], eachAction)
                
                sum = 0
                
                #Calculates the Q value of each state,action pair
                for eachSuccessor, probability in transition_states:
                    reward = self.mdp.getReward(allStates[i], eachAction, eachSuccessor)
                    sum += probability * (reward  + (discount * V_k[eachSuccessor]))
                
                Q_value_for_each_action.append(sum)
            
            #Takes care of the case where there are no transition states (i.e. a terminal state)
            if len(Q_value_for_each_action) == 0:
                V_k[allStates[i]] = 0
                policy[allStates[i]] = None
            else:
                #Value iteration takes the best Q-value (which corresponds to the best action to take)
                V_k[allStates[i]] = max(Q_value_for_each_action)
            
                #Gets the highest value's index in the Q_value list
                #best_index = np.argmax(Q_value_for_each_action)
                best_index = Q_value_for_each_action.index(max(Q_value_for_each_action))
            
                #Uses best_index to get the best action
                best_action = legal_actions[best_index]
                
                #Policy is a DICTIONARY-> For each state, there will be a corresponding best action to take.
                policy[allStates[i]] = best_action
        
            counter = counter + 1
            
            #Prevents i from going out of bounds when iterating thru all states
            i = (i + 1) % len(allStates)
            
        self.values = V_k.copy()

class PrioritizedSweepingValueIterationAgent(AsynchronousValueIterationAgent):
    """
        * Please read learningAgents.py before reading this.*

        A PrioritizedSweepingValueIterationAgent takes a Markov decision process
        (see mdp.py) on initialization and runs prioritized sweeping value iteration
        for a given number of iterations using the supplied parameters.
    """
    def __init__(self, mdp, discount = 0.9, iterations = 100, theta = 1e-5):
        """
          Your prioritized sweeping value iteration agent should take an mdp on
          construction, run the indicated number of iterations,
          and then act according to the resulting policy.
        """
        self.theta = theta
        ValueIterationAgent.__init__(self, mdp, discount, iterations)

    def runValueIteration(self):
        "*** YOUR CODE HERE ***"
                
        #Gives a list of all the states
        allStates = self.mdp.getStates()
        k = self.iterations
        discount = self.discount
        
        #Gives a list of all the states with CURRENT VALUES
        V_k = self.values
        pq = util.PriorityQueue()
        predecessors = {}
        
        #Predecessor is a dictionary that maps this: key = state and value = all parent node states (i.e. predecessors) of that (key) state
        for i in range(0, len(allStates)):
            predecessors[allStates[i]] = set()
         
        for i in range(0, len(allStates)):
            #Get all possible actions EACH STATE can take
                #Ex: Get all legal actions from state 1
            legal_actions = self.mdp.getPossibleActions(allStates[i])
            Q_value_for_each_action = []
            
            #Look at each action each state can take, and compute the Q values
            for eachAction in legal_actions:
                # In transition_states, it has a TUPLE of successor states and the probability to get to it.
                transition_states = self.mdp.getTransitionStatesAndProbs(allStates[i], eachAction)
                
                sum = 0
                
                #Calculates the Q value of each state,action pair
                for eachSuccessor, probability in transition_states:
                    #Computes the predecessors of all states:
                    #Current state we're looking at = the predecessor/parent -> values of predecessor dictionary
                    #Successor states we're looking at = keys of predecessor dictionary
                    if probability != 0:
                        predecessors[eachSuccessor].add(allStates[i])
                  
                    reward = self.mdp.getReward(allStates[i], eachAction, eachSuccessor)
                    sum += probability * (reward  + (discount * V_k[eachSuccessor]))
                
                Q_value_for_each_action.append(sum)
            
            #Takes care of the case where there are transition states (i.e. a terminal state)
            #if len(Q_value_for_each_action) != 0:
            if not ( self.mdp.isTerminal(allStates[i])) :
                #Abs value difference btwn the current value of s in self.values and highst
                #Q value across all possibile actions from s
                diff = abs(V_k[allStates[i]] - max(Q_value_for_each_action))
                #PQ contains item = (current_state, current_states_value), and priority value: -diff
                pq.push(allStates[i], -diff)
                
                
        for iteration in range(0,k):
            if pq.isEmpty():
                break
            else:
    
                #Pop a state s off the priority queue
                s = pq.pop()
                
                #Update the value of s (if its not a terminal sttae) in self.values <-> V_k
                    #Takes care of the case where there are transition states (i.e. a terminal state)
                if not ( self.mdp.isTerminal(s)):
                
                    #Calculate Q value of s:
                    legal_actions = self.mdp.getPossibleActions(s)
                    Q_value_for_each_action = []

                    #Look at each action each state can take, and compute the Q values
                    for eachAction in legal_actions:

                        # In transition_states, it has a TUPLE of successor states and the probability to get to it.
                        transition_states = self.mdp.getTransitionStatesAndProbs(s, eachAction)

                        sum = 0
                        #Calculates the Q value of each state,action pair
                        for eachSuccessor, probability in transition_states:
                            reward = self.mdp.getReward(s, eachAction, eachSuccessor)
                            sum += probability * (reward  + (discount * V_k[eachSuccessor]))

                        Q_value_for_each_action.append(sum)
                
    
                    V_k[s] = max(Q_value_for_each_action)
                    
                    for eachPredecessor in predecessors[s]:
                            
                    #Get the Q values of (predecessor,a):
                        legal_actions = self.mdp.getPossibleActions(eachPredecessor)
                        Q_value_for_each_action = []
                        
                        #Look at each action each state can take, and compute the Q values
                        for eachAction in legal_actions:
                            # In transition_states, it has a TUPLE of successor states and the probability to get to it.
                            transition_states = self.mdp.getTransitionStatesAndProbs(eachPredecessor, eachAction)
                            
                            sum = 0
                            #Calculates the Q value of each state,action pair
                            for successorState, probability in transition_states:
                                reward = self.mdp.getReward(eachPredecessor, eachAction, successorState)
                                sum += probability * (reward  + (discount * V_k[successorState]))
                            
                            Q_value_for_each_action.append(sum)
                    
                        #Calcualte the abs difference btwn the predecessor value and the highest Q value across all possible actions from p
                        diff = abs(V_k[eachPredecessor] - max(Q_value_for_each_action))
                            
                        if diff > self.theta:
                            pq.update(eachPredecessor, -diff)
